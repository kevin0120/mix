// @flow

import OWebSocket from 'ws';
import isNil from 'lodash/isNil';
import { call, take, takeLatest, put, select, fork, cancel, delay } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import type { Saga, EventChannel } from 'redux-saga';
import { onchangeIO } from '../io/action';
import { RUSH } from './action';
import { ScannerNewData } from '../scanner/action';
// import { getIBypass, getIModeSelect, handleIOFunction } from '../io/saga';
import { setHealthzCheck } from '../healthzCheck/action';
import { setNewNotification } from '../notification/action';
// import { switch2Ready, operationTrigger } from '../operation/action';
import { toolStatusChange, toolNewResults } from '../tools/action';
// import { andonScanner } from '../andon/action';
import { CommonLog } from '../../common/utils';
import type { tWebSocketEvent, tRushWebSocketData, tBarcode, tReader } from './type';
import type {tIOWSMsgType, tIOContact } from '../io/type';
import { ReaderNewData } from '../reader/action';

let task = null;
let ws = null;
const WebSocket = require('@oznu/ws-connect');


const DebounceWaitTime = 2000;

export function* watchRushEvent(): Saga<void> {
  try {
    yield takeLatest(RUSH.INIT, initRush);
    yield delay(DebounceWaitTime);
  } catch (e) {
    CommonLog.lError(e);
  }

}

function* initRush() {
  try {
    const state = yield select();

    const { connections } = state.setting.system;

    if (connections.rush === '') {
      return;
    }

    const conn = connections.rush.split('://')[1];
    const wsURL = `ws://${conn}/rush/v1/ws`;

    yield call(stopRush);

    ws = new WebSocket(wsURL, { reconnectInterval: 3000 });

    task = yield fork(
      watchRushChannel,
      state.setting.page.odooConnection.hmiSn.value
    );
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* stopRush() {
  try {
    if (isNil(ws)) {
      return;
    }
    if (
      ws.ws.readyState === OWebSocket.OPEN ||
      ws.ws.readyState === OWebSocket.CONNECTING
    ) {
      yield put(setHealthzCheck('masterpc', false));
      yield put(setNewNotification('Info', `masterPC连接状态更新: ${false}`));
      yield put(setHealthzCheck('controller', false));
      yield put(setNewNotification('Info', `controller连接状态更新: ${false}`));
      ws.close();
    }
    ws = null;
    if (task) {
      yield cancel(task);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function createRushChannel(hmiSN: string): EventChannel<void> {
  return eventChannel(emit => {
    ws.on('open', () => {
      emit({ type: 'healthz', payload: true });
      // reg msg
      ws.sendJson({ hmi_sn: hmiSN }, err => {
        if (err) {
          ws.close();
        }
      });
    });

    ws.on('close', () => {
      emit({ type: 'healthz', payload: false });
    });

    ws.on('error', () => {
      emit({ type: 'healthz', payload: false });
      // console.log('websocket error. reconnect after 1s');
    });
    ws.on('ping', () => {
      CommonLog.Debug('receive ping msg');
    });
    ws.on('pong', () => {
      CommonLog.Debug('receive pong msg');
    });

    ws.on('message', data => {
      emit({ type: 'data', payload: data });
    });

    return () => {
    };
  });
}

function* handleRushData(type: tWebSocketEvent, data: tRushWebSocketData): Saga<void> {
  try {
    switch (type) {
      case 'maintenance':
        yield put(setNewNotification('Maintenance', `新维护请求: ${data.type},${data.data.name}`));
        break;
      case 'job':
        // CommonLog.Info(json);
        // if (state.workMode.workMode === 'manual' && json.job_id > 0) {
        //   if (state.setting.operationSettings.manualFreestyle) {
        //     yield put(switch2Ready(false));
        //     // state=yield select();
        //   }
        //   const { carID } = state.operations;
        //
        //   yield put(operationTrigger(
        //     carID,
        //     '',
        //     json.job_id,
        //     OPERATION_SOURCE.MANUAL
        //   ));
        // }
        break;
      case 'odoo': {
        // const odooHealthzStatus = json.status === 'online';
        // const healthzStatus = state.healthCheckResults; // 获取整个healthz
        // if (
        //   !lodash.isEqual(healthzStatus.odoo.isHealth, odooHealthzStatus)
        // ) {
        //   yield put(setHealthzCheck('odoo', odooHealthzStatus));
        //   yield put(
        //     setNewNotification(
        //       'info',
        //       `后台连接状态更新: ${odooHealthzStatus}`
        //     )
        //   );
        // }
        break;
      }
      case 'io': {
        let d;
        const msgType = (data.type: tIOWSMsgType);
        switch (msgType) {
          case 'WS_IO_CONTACT': {
            d = (data.data: tIOContact);
            break;
          }
          default:
            CommonLog.lError('IO Message Type Is Not Defined');
        }
        yield put(onchangeIO(d));
        break;
      }
      case 'result':
        CommonLog.Info(` tool new results: ${data.data}`);
        yield put(toolNewResults(data.data));
        break;
      case 'scanner': {
        const d = (data.data: tBarcode);
        CommonLog.Info(` Scanner receive data: ${d.barcode}`);
        yield put(ScannerNewData(d.barcode));
        break;
      }
      case 'reader': {
        const d = (data.data: tReader);
        CommonLog.Info(` Reader receive data: ${d.uid}`);
        yield put(ReaderNewData(d.uid));
        break;
      }

      case 'controller': {
        // const healthzStatus = state.healthCheckResults; // 获取整个healthz
        // const controllerHealthzStatus = data.status === 'online';
        // if (
        //   !lodash.isEqual(
        //     healthzStatus.controller.isHealth,
        //     controllerHealthzStatus
        //   )
        // ) {
        //   // 如果不相等 更新
        //   yield put(
        //     setHealthzCheck('controller', controllerHealthzStatus)
        //   );
        //   yield put(
        //     setNewNotification(
        //       'info',
        //       `controller连接状态更新: ${controllerHealthzStatus}`
        //     )
        //   );
        // }
        break;
      }

      case 'tool': {
        yield put(toolStatusChange(data.tool_sn, data.status, data.reason));
        break;
      }

      case 'tightening_device': {
        // 初始化所有拧紧设备
        break;
      }

      default:
        break;
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

export function* watchRushChannel(hmiSN: string): Saga<void> {
  try {
    const chan = yield call(createRushChannel, hmiSN);
    while (true) {
      const data = yield take(chan);
      // let state = yield select();

      const { type, payload } = data;

      switch (type) {
        case 'healthz': {
          // const healthzStatus = state.healthCheckResults; // 获取整个healthz
          // if (!lodash.isEqual(healthzStatus.masterpc.isHealth, payload)) {
          //   // 如果不相等 更新
          //   yield put(setHealthzCheck('masterpc', payload));
          //   yield put(
          //     setNewNotification('info', `masterPC连接状态更新: ${payload}`)
          //   );
          // }
          // if (!payload) {
          //   yield put(setHealthzCheck('controller', false));
          //   yield put(
          //     setNewNotification('info', `controller连接状态更新: ${false}`)
          //   );
          // }
          break;
        }
        case 'data': {
          const dataArray = payload.split(';');
          const event: tWebSocketEvent = dataArray[0].split(':').slice(-1)[0];

          const json: tRushWebSocketData = JSON.parse(dataArray.slice(-1));

          yield fork(handleRushData, event, json); // 异步处理rush数据
          break;
        }
        default:
          break;
      }
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}
