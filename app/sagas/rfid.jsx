import {
  call,
  take,
  put,
  select,
  fork,
  cancel,
  // debounce,
  delay
} from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import { RFID } from '../actions/actionTypes';
import { triggerOperation } from './operation';
import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';
import { setHealthzCheck } from '../actions/healthCheck';
import { setNewNotification } from '../actions/notification';

const _ = require('lodash');

let client = null;
let recon = null;

let watchChannelTask = null;

const net = require('net');
const Reconnect = require('node-net-reconnect');

const lodash = require('lodash');

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
  const { connections ,setting:{operationSettings}} = state;

  const { regExp } = operationSettings;

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

  const rfidChannel = yield call(createRfidChannel, client);

  watchChannelTask = yield fork(watchRfidChannel, rfidChannel, regExp);
}

export function* stopRFID() {
  try {
    if (recon != null) {
      recon.end();
      recon = null;
    }

    if (client) {
      client.destroy();
      client = null;
    }
    if (watchChannelTask){
      yield cancel(watchChannelTask);
    }
  }catch (e) {
    console.error(`stopRFID error: ${e.message}`)
  }

}

function createRfidChannel(rfidClient) {
  return eventChannel(emit => {
    rfidClient.on('connect', () => {
      // health
      emit({ type: 'healthz', payload: true });
    });

    rfidClient.on('data', data => {
      emit({ type: 'data', payload: data });
    });

    rfidClient.on('end', () => {
      emit({ type: 'healthz', payload: false });
    });

    rfidClient.on('close', () => {
      // unhealth
      emit({ type: 'healthz', payload: false });
      rfidClient.end(); // try to reconnect
    });

    rfidClient.on('error', () => {
      // unhealth
      emit({ type: 'healthz', payload: false });
      rfidClient.end(); // try to reconnect
    });

    rfidClient.on('timeout', () => {
      // timeout
    });

    return () => {};
  });
}

// const getHealthz = state => state.healthCheckResults;

function* RFIDHandler(reg, data) {
  try {

    const state = yield select();

    const { rfidEnabled } = state.setting.systemSettings;

    if (!rfidEnabled) {
      // 未使能rfid
      return
    }
    const { type, payload } = data;

    switch (type) {
      case 'healthz': {
        const {healthCheckResults: healthzStatus} = state; // 获取整个healthz
        if (!lodash.isEqual(healthzStatus.rfid.isHealth, payload)) {
          // 如果不相等 更新
          yield put(setHealthzCheck('rfid', payload));
          yield put(setNewNotification('info', `rfid连接状态更新: ${payload}`));
        }
        break;
      }
      case 'data': {
        if (payload !== 'NoRead') {
          const buf2 = Buffer.from(payload, 'hex');
          const code = buf2.toString();
          const dataValue = _.trim(code);

          if (!reg.test(dataValue)) {
            yield put(
              setNewNotification(
                'error',
                `RFID data can not match: ${dataValue}`
              )
            );
          }
          const targetValue = reg.exec(dataValue)[0];
          // yield put(setNewNotification('info', `new rfid value: ${targetValue}`));

          if (isCarID(targetValue)) {
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
        break;
      }
      default:
        break;
    }
  } catch (err) {
    console.log(`rfid error msg:${err.message}`);
  }
}

export function* watchRfidChannel(channel,regExp) {
  const reg = new RegExp(regExp, 'i'); // 正則表達式,大小寫不敏感
  while (client !== null) {
    try {
      // yield debounce(3000, rfidChannel, RFIDHandler, reg); // RFID 因为频繁触发所以进行防抖动处理,默认3秒
      const data =yield take(channel);
      yield fork(RFIDHandler,data,reg);
      yield delay(2000);
    } catch (err) {
      console.error(`watchRFIDChannel: ${err.message}`);
    }
  }
}
