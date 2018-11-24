import { call, take, select, put, fork } from 'redux-saga/effects';
import OWebSocket from "ws";
import { eventChannel } from 'redux-saga';
import { ANDON, AIIS } from '../actions/actionTypes';
import { jobManual } from './api/operation';
import { setHealthzCheck } from '../actions/healthCheck';

let ws = null;

const WebSocket = require('@oznu/ws-connect');

const AIIS_WS_CHANNEL={
  OPEN:'AIIS_WS_CHANNEL_OPEN',
  CLOSE:'AIIS_WS_CHANNEL_CLOSE',
  ERROR:'AIIS_WS_CHANNEL_ERROR',
  PING:'AIIS_WS_CHANNEL_PING',
  PONG:'AIIS_WS_CHANNEL_PONG',
  MESSAGE:'AIIS_WS_CHANNEL_PONG'
};

export function* watchAiis() {
  while (true) {
    const action = yield take([ANDON.NEW_DATA, AIIS.INIT]);
    switch (action.type) {
      case ANDON.NEW_DATA:
        yield call(handleAiisData, action.json);
        break;
      case AIIS.INIT:
        yield fork(initAiis,action.aiisUrl,action.hmiSN);
        break;
      default:
        break;
    }
  }
}

function* initAiis(aiisUrl,hmiSN){
  if (ws) {
    yield call(stopAiisWebsocket);
  }
  ws = new WebSocket(aiisUrl);

  yield fork(aiisWSListener,hmiSN);
  yield call(wsOnOpen, hmiSN);
}

export function* handleAiisData(data) {
  const state = yield select();
  if (state.operations.operationStatus !== OPERATION_STATUS.DOING) {
    if (data.vin_code.length) {
      // 车辆拧紧作业
      yield call(triggerOperation, data.vin_code, data.cartype_code, null, OPERATION_SOURCE.ANDON);
    } else {
      // 空车信息

      const {
        carType,
        carID,
        productID,
        workcenterID,
        results
      } = state.operations;

      const { jobID } = state.setting.operationSettings.emptyCarJob;

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
        jobID,
        results,
        hmiSn.value,
        productID,
        workcenterID,
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
        type: AIIS_WS_CHANNEL.CLOSE ,
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

    ws.on('message', (dataRaw) => {
      emit({
        type: AIIS_WS_CHANNEL.PONG ,
        dataRaw
      });
    });

    return () => {
    };
  });
}

function* aiisWSListener(hmiSN) {
  const chan = yield call(aiisWSChannel);
  try {
    while (true) {
      const chanAction = yield take(chan);
      switch (chanAction.type) {
        case AIIS_WS_CHANNEL.OPEN:
          yield call(wsOnOpen,hmiSN);
          break;
        case AIIS_WS_CHANNEL.CLOSE:
          yield put(setHealthzCheck('Andon', false));
          console.log(
            `websocket disconnected. retry in 1s code: ${chanAction.code}, reason: ${chanAction.reason}`
          );
          break;
        case AIIS_WS_CHANNEL.ERROR:
          yield put(setHealthzCheck('Andon', false));
          console.log('websocket error. reconnect after 1s');
          break;
        case AIIS_WS_CHANNEL.PING:
          console.log(' receive ping Msg');
          break;
        case AIIS_WS_CHANNEL.PONG:
          console.log(' receive pong Msg');
          break;
        case AIIS_WS_CHANNEL.MESSAGE:
          yield call(wsOnMessage,chanAction.dataRaw);
          break;
        default:
          break;
      }
    }
  }
  finally {
    console.log('aiisWSListener finished');
  }
}

function* wsOnOpen(hmiSN){
  ws.sendJson({ hmi_sn: hmiSN }, err => {
    if (err) {
      console.log('aiis ws sendJson error');
      ws.close();
    }
  });
  yield put(setHealthzCheck('andon', true));
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


