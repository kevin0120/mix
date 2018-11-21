import OWebSocket from 'ws';
import { call, take, put, select, fork } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga'
import { RUSH } from '../actions/actionTypes';
import { NewResults } from '../actions/rush';

let rushWS = null;
let rushChannel = null;
const WebSocket = require('@oznu/ws-connect');

export function* watchRush() {
  while (true) {
    const { type } = yield take(RUSH.INIT);
    switch (type) {
      case RUSH.INIT:
        yield call(initRush);
        break;

      default:
        break;
    }
  }
}

function* initRush() {
  const state = yield select();

  const {
    connections,
    setting
  } = state;

  if (connections.masterpc === '') {
    return;
  }

  const conn = connections.masterpc.split('://')[1];
  const wsURL = `ws://${conn}/rush/v1/ws`;

  if (rushWS) {
    yield call(stopRush);
  }

  rushWS = new WebSocket(wsURL);
  rushChannel = yield call(createRushChannel, rushWS, setting.page.odooConnection.hmiSn.value);

  yield fork(watchRushChannel);
}

function stopRush() {
  if (
    rushWS.ws.readyState === OWebSocket.OPEN ||
    rushWS.ws.readyState === OWebSocket.CONNECTING
  ) {
    rushWS.close();
  }
  rushWS = null;
}

function createRushChannel(ws, hmiSN) {
  return eventChannel(emit => {
    ws.on('open', () => {
      // reg msg
      ws.sendJson({ hmi_sn: hmiSN }, err => {
        if (err) {
          ws.close();
        }
      });
    });

    ws.on('close', () => {

    });

    ws.on('error', () => {
      // console.log('websocket error. reconnect after 1s');
    });
    ws.on('ping', () => {
      // console.log(' receive ping Msg');
    });
    ws.on('pong', () => {
      // console.log(' receive pong Msg');
    });

    ws.on('message', (dataRaw) => {
      emit(dataRaw);
    });

    return () => {
    };
  });
}

export function* watchRushChannel() {
  while (true) {
    const payload = yield take(rushChannel);

    const dataArray = payload.split(';');
    const event = dataArray[0].split(':').slice(-1)[0];

    const data = dataArray.slice(-1);
    const json = JSON.parse(data);

    switch (event) {
      case 'job':
        break;
      case 'io':
        break;
      case 'result':
        yield put(NewResults(json));
        break;
      case 'scanner':
        break;
      default:
        break;
    }
  }
}
