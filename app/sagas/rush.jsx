import OWebSocket from 'ws';
import {
  call,
  take,
  takeLatest,
  put,
  select,
  fork,
  cancel
} from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import { OPERATION, RUSH, WORK_MODE } from '../actions/actionTypes';
import { NewResults } from '../actions/rush';
import { NewCar } from '../actions/scannerDevice';
import { getIBypass, getIModeSelect, handleIOFunction } from './io';
import { OPERATION_SOURCE } from '../reducers/operations';
import { IO_FUNCTION } from '../reducers/io';
import { setHealthzCheck } from '../actions/healthCheck';
import { setNewNotification } from '../actions/notification';
import { switch2Ready, operationTrigger } from '../actions/operation';

let task = null;
let ws = null;
const WebSocket = require('@oznu/ws-connect');

const lodash = require('lodash');

export function* watchRush() {
  yield takeLatest(RUSH.INIT, initRush);
}

function* initRush() {
  try {
    const state = yield select();

    const { connections } = state;

    if (connections.masterpc === '') {
      return;
    }

    const conn = connections.masterpc.split('://')[1];
    const wsURL = `ws://${conn}/rush/v1/ws`;

    yield call(stopRush);

    ws = new WebSocket(wsURL, { reconnectInterval: 3000 });

    task = yield fork(
      watchRushChannel,
      state.setting.page.odooConnection.hmiSn.value
    );
  } catch (e) {
    console.log(e);
  }
}

function* stopRush() {
  try {
    if (ws) {
      if (
        ws.ws.readyState === OWebSocket.OPEN ||
        ws.ws.readyState === OWebSocket.CONNECTING
      ) {
        yield put(setHealthzCheck('masterpc', false));
        yield put(setNewNotification('info', `masterPC连接状态更新: ${false}`));
        yield put(setHealthzCheck('controller', false));
        yield put(setNewNotification('info', `controller连接状态更新: ${false}`));
        ws.close();
      }
      ws = null;
    }
    if (task) {
      yield cancel(task);
    }
  } catch (e) {
    console.error(`stopRush error: ${e.message}`);
  }


}

function createRushChannel(hmiSN) {
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
      // console.log(' receive ping Msg');
    });
    ws.on('pong', () => {
      // console.log(' receive pong Msg');
    });

    ws.on('message', data => {
      emit({ type: 'data', payload: data });
    });

    return () => {
    };
  });
}

// const getHealthz = state => state.healthCheckResults;

export function* watchRushChannel(hmiSN) {
  try {
    const chan = yield call(createRushChannel, hmiSN);
    while (true) {
      const data = yield take(chan);
      const state = yield select();

      const { type, payload } = data;

      switch (type) {
        case 'healthz': {
          const healthzStatus = state.healthCheckResults; // 获取整个healthz
          if (!lodash.isEqual(healthzStatus.masterpc.isHealth, payload)) {
            // 如果不相等 更新
            yield put(setHealthzCheck('masterpc', payload));
            yield put(
              setNewNotification('info', `masterPC连接状态更新: ${payload}`)
            );
          }
          if (!payload) {
            yield put(setHealthzCheck('controller', false));
            yield put(
              setNewNotification('info', `controller连接状态更新: ${false}`)
            );
          }
          break;
        }
        case 'data': {
          const dataArray = payload.split(';');
          const event = dataArray[0].split(':').slice(-1)[0];

          const json = JSON.parse(dataArray.slice(-1));

          switch (event) {
            case 'maintenance':
              yield put(
                setNewNotification(
                  'maintenance',
                  `新维护请求: ${json.type},${json.name}`
                )
              );
              break;
            case 'job':
              if (state.workMode.workMode === 'manual' && json.job_id > 0) {
                if (state.setting.operationSettings.manualFreestyle) {
                  yield put(switch2Ready());
                }

                const { carID, carType } = state.operations;
                yield put(operationTrigger(
                  carID,
                  carType,
                  json.job_id,
                  OPERATION_SOURCE.MANUAL
                ));
              }

              break;
            case 'odoo': {
              const odooHealthzStatus = json.status === 'online';
              const healthzStatus = state.healthCheckResults; // 获取整个healthz
              if (
                !lodash.isEqual(healthzStatus.odoo.isHealth, odooHealthzStatus)
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
            case 'scanner': {
              // console.log('rush scanner:', json);
              const { workMode } = state.workMode;
              if (workMode === 'scanner' || workMode === 'manual') {
                yield put(NewCar(json.barcode));
              }
              break;
            }
            case 'controller': {
              const healthzStatus = state.healthCheckResults; // 获取整个healthz
              const controllerHealthzStatus = json.status === 'online';
              if (
                !lodash.isEqual(
                  healthzStatus.controller.isHealth,
                  controllerHealthzStatus
                )
              ) {
                // 如果不相等 更新
                yield put(
                  setHealthzCheck('controller', controllerHealthzStatus)
                );
                yield put(
                  setNewNotification(
                    'info',
                    `controller连接状态更新: ${controllerHealthzStatus}`
                  )
                );
              }
              break;
            }

            case 'tool': {
              console.log(json);
              break;
            }
            default:
              break;
          }
          break;
        }
        default:
          break;
      }
    }
  } catch (e) {
    console.log(e);
  }
}
