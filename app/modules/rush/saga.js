// @flow

import OWebSocket from 'ws';
import { call, put, select, fork, cancel, takeEvery } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import type { Saga, EventChannel } from 'redux-saga';
import { RUSH } from './action';
import { setHealthzCheck } from '../healthzCheck/action';
import { setNewNotification } from '../notification/action';
import { CommonLog } from '../../common/utils';
import handleData from './handleData';
import handleHealthz from './handleHealthz';
import { WEBSOCKET_EVENTS } from './type';

let task = null;
let ws = null;
const WebSocket = require('@oznu/ws-connect');


export function getWSClient() {
  return ws;
}

export function* watchRushEvent(): Saga<void> {
  try {
    yield takeEvery(RUSH.INIT, initRush);
  } catch (e) {
    CommonLog.lError(e, { at: 'watchRushEvent' });
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
    CommonLog.lError(e, { at: 'initRush' });
  } finally {
    if (!(ws && task)) {
      if (ws) {
        ws.close();
        ws = null;
      }
      if (task) {
        yield cancel(task);
        task = null;
      }
    }
  }
}

function* stopRush() {
  try {
    if (task) {
      yield cancel(task);
      task = null;
    }
    if (!ws) {
      return;
    }
    if (
      ws.ws.readyState === OWebSocket.OPEN ||
      ws.ws.readyState === OWebSocket.CONNECTING
    ) {
      yield put(setHealthzCheck('masterpc', false));
      yield put(setNewNotification('Info', `masterPC连接状态更新: false`));
      yield put(setHealthzCheck('controller', false));
      yield put(setNewNotification('Info', `controller连接状态更新: false`));
      ws.close();
    }
    ws = null;
  } catch (e) {
    CommonLog.lError(e, { at: 'stopRush' });
  }
}

function createRushChannel(hmiSN: string): EventChannel<void> {
  return eventChannel(emit => {
    if (ws) {
      ws.on('open', () => {
        emit({ type: 'healthz', payload: true });
        // reg msg
        if (ws) {
          ws.sendJson({
            type:'WS_REG',
            sn:0,
            data:{
              hmi_sn: hmiSN
            }
          }, err => {
            if (err && ws) {
              CommonLog.lError(err);
              ws.close();
            }
          });
        }
      });

      ws.on('close', (...args) => {
        console.warn('close', ...args, ws);

        emit({ type: 'healthz', payload: false });
      });

      ws.on('error', (...args) => {
        console.warn('error', ...args);

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
        if (data.type === WEBSOCKET_EVENTS.reply) {
          emit({ type: 'reply', payload: data });
          return;
        }
        emit({ type: 'data', payload: data });
      });
      ws.on('websocket-status', (...args) => {

      })
    } else {
      CommonLog.lError('ws doesn\'t exist', { at: 'createRushChannel' });
    }
    return () => {
    };
  });
}

const rushChannelHandlers = {
  healthz: handleHealthz,
  data: handleData
};

function* rushChannelTask(data) {
  try {
    yield call(rushChannelHandlers[data.type], data.payload);
  } catch (e) {
    CommonLog.lError(e, { at: 'rushChannelTask' });
  }
}

export function* watchRushChannel(hmiSN: string): Saga<void> {
  try {
    const chan = yield call(createRushChannel, hmiSN);
    yield takeEvery(chan, rushChannelTask);
  } catch (e) {
    CommonLog.lError(e, { at: 'watchRushChannel' });
  }
}
