import { call, take, put, select, fork } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import { RFID } from '../actions/actionTypes';
import { triggerOperation } from './operation';
import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';

const _ = require('lodash');

let client = null;
let recon = null;
let rfidChannel = null;

const net = require('net');
const Reconnect = require('node-net-reconnect');

export function* watchRfid() {
  while (true) {
    const { type } = yield take(RFID.INIT);
    switch (type) {
      case RFID.INIT:
        yield call(initRFID);
        break;

      default:
        break;
    }
  }
}

function* initRFID() {
  const state = yield select();

  const { connections } = state;

  if (connections.rfid === '') {
    return;
  }

  const rfidURL = connections.rfid;
  const kvs = rfidURL.split('://');
  const hostPorts = kvs[1].split(':');
  const host = hostPorts[0];
  const port = parseInt(hostPorts[1].split('/')[0], 10);

  yield call(stopRFID);

  client = new net.Socket();
  const options = {
    host,
    port,
    retryTime: 1000, // 1s for every retry
    retryAlways: true // retry even if the connection was closed on purpose
  };

  recon = new Reconnect(client, options);
  client.setTimeout(10000);
  client.setEncoding('ascii');
  client.connect(options);

  rfidChannel = yield call(createRfidChannel, client);

  yield fork(watchRfidChannel);
}

export function stopRFID() {
  if (recon != null) {
    recon.end();
    recon = null;
  }

  if (client) {
    client.destroy();
    client = null;
  }
}

function createRfidChannel(rfidClient) {
  return eventChannel(emit => {
    rfidClient.on('connect', () => {
      // health
    });

    rfidClient.on('data', data => {
      emit(data);
    });

    rfidClient.on('end', () => {
      // unhealth
    });

    rfidClient.on('close', () => {
      // unhealth
      rfidClient.end(); // try to reconnect
    });

    rfidClient.on('error', () => {
      // unhealth
      rfidClient.end(); // try to reconnect
    });

    rfidClient.on('timeout', () => {
      // timeout
    });

    return () => {};
  });
}

export function* watchRfidChannel() {
  while (client !== null) {
    const payload = yield take(rfidChannel);

    if (payload !== 'NoRead') {
      const buf2 = Buffer.from(_.trim(payload), 'hex');
      const code = buf2.toString();
      const dataValue = _.trim(code);

      if (isCarID(dataValue)) {
        yield call(
          triggerOperation,
          dataValue,
          null,
          null,
          OPERATION_SOURCE.RFID
        );
      } else {
        yield call(
          triggerOperation,
          null,
          dataValue,
          null,
          OPERATION_SOURCE.RFID
        );
      }
    }
  }
}
