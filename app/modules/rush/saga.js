// @flow

import OWebSocket from 'ws';
import { call, take, takeLatest, put, select, fork, cancel, delay, takeEvery } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import type { Saga, EventChannel } from 'redux-saga';
import { RUSH } from './action';
import { setHealthzCheck } from '../healthzCheck/action';
import { setNewNotification } from '../notification/action';
import { CommonLog } from '../../common/utils';
import handleData from './handleData';
import handleHealthz from './handleHealthz';
import { WEBSOCKET_EVENTS } from './type';

const tasks = [];
let ws = null;

const WebSocket = require('@oznu/ws-connect');

const DebounceWaitTime = 2000;

export function getWSClient() {
  return ws;
}

export function* watchRushEvent(): Saga<void> {
  try {
    yield takeLatest(RUSH.INIT, initRush);
    yield delay(DebounceWaitTime);
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

    const listenerTask = yield fork(
      watchRushChannel,
      state.setting.page.odooConnection.hmiSn.value
    );
    tasks.push(listenerTask);


  } catch (e) {
    CommonLog.lError(e, { at: 'initRush' });
  }
}

function* stopRush() {
  try {
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
    while (tasks.length > 0) {
      yield cancel(tasks.pop());
    }
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
          ws.sendJson({ hmi_sn: hmiSN }, err => {
            if (err && ws) {
              ws.close();
            }
          });
        }
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
        if (data.type === WEBSOCKET_EVENTS.reply) {
          emit({ type: 'reply', payload: data });
          return;
        }
        emit({ type: 'data', payload: data });
      });
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
