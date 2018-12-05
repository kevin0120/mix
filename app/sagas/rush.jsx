import OWebSocket from 'ws';
import { call, take, put, select, fork } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import { OPERATION, RUSH, WORK_MODE } from '../actions/actionTypes';
import { NewResults } from '../actions/rush';
import { NewCar } from '../actions/scannerDevice';
import { getIBypass, getIModeSelect, handleIOFunction } from './io';
import { triggerOperation } from './operation';
import { OPERATION_SOURCE } from '../reducers/operations';
import { IO_FUNCTION } from '../reducers/io';
import { setHealthzCheck } from "../actions/healthCheck";
import { setNewNotification } from "../actions/notification";

let rushWS = null;
let rushChannel = null;
const WebSocket = require('@oznu/ws-connect');

const lodash = require('lodash');

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

  const { connections, setting } = state;

  if (connections.masterpc === '') {
    return;
  }

  const conn = connections.masterpc.split('://')[1];
  const wsURL = `ws://${conn}/rush/v1/ws`;

  if (rushWS) {
    yield call(stopRush);
  }

  rushWS = new WebSocket(wsURL);
  rushChannel = yield call(
    createRushChannel,
    rushWS,
    setting.page.odooConnection.hmiSn.value
  );

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

    ws.on('close', () => {});

    ws.on('error', () => {
      // console.log('websocket error. reconnect after 1s');
    });
    ws.on('ping', () => {
      // console.log(' receive ping Msg');
    });
    ws.on('pong', () => {
      // console.log(' receive pong Msg');
    });

    ws.on('message', dataRaw => {
      emit(dataRaw);
    });

    return () => {};
  });
}

const getHealthz = state => state.healthCheckResults;


export function* watchRushChannel() {
  while (rushWS !== null) {
    const payload = yield take(rushChannel);

    const dataArray = payload.split(';');
    const event = dataArray[0].split(':').slice(-1)[0];

    const data = dataArray.slice(-1);
    const json = JSON.parse(data);

    const state = yield select();
    switch (event) {
      case 'job':
        if (state.workMode.workMode === 'manual' && json.job_id) {
          yield call(
            triggerOperation,
            null,
            null,
            json.job_id,
            OPERATION_SOURCE.MANUAL
          );
        }

        break;
      case 'odoo': {
        const odooHealthzStatus = json.status === 'online';
        const healthzStatus = yield select(getHealthz); // 获取整个healthz
        if (
          !lodash.isEqual(
            healthzStatus.odoo.isHealth,
            odooHealthzStatus
          )
        ) {
          yield put(setHealthzCheck('odoo', odooHealthzStatus));
          yield put(
            setNewNotification(
              'info',
              `后台连接状态更新: ${odooHealthzStatus}`
            )
          );
        }
        break;
      }
      case 'io':
        if (state.setting.systemSettings.modbusEnable) {
          break;
        }

        if (json.inputs) {
          if (json.inputs[getIBypass()] === '1') {
            // 强制放行
            // yield put({ type: OPERATION.FINISHED, data: [] });
            yield call(handleIOFunction, IO_FUNCTION.IN.BYPASS);
          }

          if (json.inputs[getIModeSelect()] === '1') {
            // 切换到手动模式
            yield put({ type: WORK_MODE.SWITCH_WM, mode: 'manual' });
          } else {
            // 切换到自动模式
            yield put({ type: WORK_MODE.SWITCH_WM, mode: 'auto' });
          }
        }
        break;
      case 'result':
        yield put(NewResults(json));
        break;
      case 'scanner':
        yield put(NewCar(json.barcode));
        break;
      default:
        break;
    }
  }
}
