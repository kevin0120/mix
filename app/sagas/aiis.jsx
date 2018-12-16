import {
  call,
  take,
  takeLatest,
  takeEvery,
  select,
  put,
  fork,
  cancel
} from 'redux-saga/effects';
import OWebSocket from 'ws';
import { eventChannel } from 'redux-saga';
import { ANDON, AIIS, OPERATION } from '../actions/actionTypes';
import { jobManual } from './api/operation';
import { setHealthzCheck } from '../actions/healthCheck';

import { triggerOperation } from './operation';

import { OPERATION_STATUS, OPERATION_SOURCE } from '../reducers/operations';

let task = null;

let ws = null;

const WebSocket = require('@oznu/ws-connect');

const AIIS_WS_CHANNEL = {
  OPEN: 'AIIS_WS_CHANNEL_OPEN',
  CLOSE: 'AIIS_WS_CHANNEL_CLOSE',
  ERROR: 'AIIS_WS_CHANNEL_ERROR',
  PING: 'AIIS_WS_CHANNEL_PING',
  PONG: 'AIIS_WS_CHANNEL_PONG',
  MESSAGE: 'AIIS_WS_CHANNEL_MESSAGE'
};

export function* watchAiis() {
  yield takeLatest(AIIS.INIT, initAiis);
  yield takeEvery(ANDON.NEW_DATA, handleAiisData);
}

function* initAiis() {
  try {
    const state = yield select();
    const aiisUrl = state.setting.system.connections.aiis;
    const hmiSN = state.setting.page.odooConnection.hmiSn.value;

    if (task) {
      yield cancel(task);
    }

    if (ws) {
      yield call(stopAiisWebsocket);
    }

    const uris = aiisUrl.split('://');
    if (uris.length > 1) {
      const url = `ws://${uris[1]}/aiis/v1/ws`;
      ws = new WebSocket(url, { reconnectInterval: 3000 });

      task = yield fork(aiisWSListener, hmiSN);
      // yield call(wsOnOpen, hmiSN);
    }
  } catch (e) {
    console.log(e);
  }
}

export function* handleAiisData(action) {
  const data = action.json;
  const state = yield select();
  if (state.operations.operationStatus !== OPERATION_STATUS.DOING) {
    if (data.vin_code.length) {
      // 车辆拧紧作业
      yield call(
        triggerOperation,
        data.vin_code,
        data.cartype_code,
        null,
        OPERATION_SOURCE.ANDON
      );
    } else {
      // 空车信息

      const { carType, carID } = state.operations;

      const { emptyCarJob } = state.setting.operationSettings;

      const controllerSN = state.connections.controllers[0].serial_no;
      const rushUrl = state.connections.masterpc;
      const { hmiSn } = state.setting.page.odooConnection;
      const userID = 1;
      const skip = true;
      const hasSet = false;

      const resp = yield call(
        jobManual,
        rushUrl,
        controllerSN,
        carType,
        carID,
        userID,
        emptyCarJob,
        hmiSn.value,
        0,
        skip,
        hasSet,
        ''
      );

      if (resp.statusCode !== 200) {
        yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
      }
    }
  }
}

function aiisWSChannel() {
  return eventChannel(emit => {
    ws.on('open', () => {
      emit({ type: AIIS_WS_CHANNEL.OPEN });
    });

    ws.on('close', (code, reason) => {
      emit({
        type: AIIS_WS_CHANNEL.CLOSE,
        code,
        reason
      });
    });

    ws.on('error', () => {
      emit({ type: AIIS_WS_CHANNEL.ERROR });
    });

    ws.on('ping', () => {
      emit({ type: AIIS_WS_CHANNEL.PING });
    });

    ws.on('pong', () => {
      emit({ type: AIIS_WS_CHANNEL.PONG });
    });

    ws.on('message', dataRaw => {
      emit({
        type: AIIS_WS_CHANNEL.MESSAGE,
        dataRaw
      });
    });

    return () => {};
  });
}

function* aiisWSListener(hmiSN) {
  const chan = yield call(aiisWSChannel);
  try {
    while (true) {
      const chanAction = yield take(chan);
      switch (chanAction.type) {
        case AIIS_WS_CHANNEL.OPEN:
          yield put(setHealthzCheck('andon', true));
          yield call(wsOnOpen, hmiSN);
          break;
        case AIIS_WS_CHANNEL.CLOSE:
          yield put(setHealthzCheck('andon', false));
          console.log(
            `websocket disconnected. retry in 1s code: ${
              chanAction.code
            }, reason: ${chanAction.reason}`
          );
          break;
        case AIIS_WS_CHANNEL.ERROR:
          yield put(setHealthzCheck('andon', false));
          console.log('websocket error. reconnect after 1s');
          break;
        case AIIS_WS_CHANNEL.PING:
          console.log(' receive ping Msg');
          break;
        case AIIS_WS_CHANNEL.PONG:
          console.log(' receive pong Msg');
          break;
        case AIIS_WS_CHANNEL.MESSAGE:
          yield call(wsOnMessage, chanAction.dataRaw);
          break;
        default:
          break;
      }
    }
  } finally {
    console.log('aiisWSListener finished');
  }
}

function* wsOnOpen(hmiSN) {
  ws.sendJson({ hmi_sn: hmiSN }, err => {
    if (err) {
      console.log('aiis ws sendJson error');
      ws.close();
    }
  });
}

function* wsOnMessage(dataRaw) {
  const dataArray = dataRaw.split(';');

  const data = dataArray.slice(-1);
  const json = JSON.parse(data);

  yield put({ type: ANDON.NEW_DATA, json });
}

function stopAiisWebsocket() {
  if (
    ws.ws.readyState === OWebSocket.OPEN ||
    ws.ws.readyState === OWebSocket.CONNECTING
  ) {
    ws.close();
  }
  ws = null;
}
