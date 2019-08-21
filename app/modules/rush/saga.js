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
import rushHealthz from './rushHealthz';
import {getWSClient, setWSClient} from './client';

let task = null;
const WebSocket = require('@oznu/ws-connect');



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

    setWSClient(new WebSocket(wsURL,
      {
        reconnectInterval: 3000,
        options:
          {
            maxPayload: 200 * 1024 * 1024
          }
      }));

    task = yield fork(
      watchRushChannel,
      state.setting.page.odooConnection.hmiSn.value
    );

  } catch (e) {
    CommonLog.lError(e, { at: 'initRush' });
  } finally {
    // const ws=getWSClient();
    // if (!(ws && task)) {
    //   if (ws) {
    //     ws.close();
    //     setWSClient(null);
    //   }
    //   if (task) {
    //     yield cancel(task);
    //     task = null;
    //   }
    // }
  }
}

function* stopRush() {
  try {
    if (task) {
      yield cancel(task);
      task = null;
    }
    if (!getWSClient()) {
      return;
    }
    if (
      getWSClient().ws.readyState === OWebSocket.OPEN ||
      getWSClient().ws.readyState === OWebSocket.CONNECTING
    ) {
      yield put(setHealthzCheck('masterpc', false));
      yield put(setNewNotification('Info', `masterPC连接状态更新: false`));
      yield put(setHealthzCheck('controller', false));
      yield put(setNewNotification('Info', `controller连接状态更新: false`));
      getWSClient().close();
    }
    setWSClient(null);
  } catch (e) {
    CommonLog.lError(e, { at: 'stopRush' });
  }
}

function createRushChannel(hmiSN: string): EventChannel<void> {
  return eventChannel(emit => {
    const ws=getWSClient();
    if (ws) {
      ws.on('open', () => {
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
              // ws.close();
            }
          });
        }
        emit({ type: 'healthz', payload: true });
      });

      ws.on('close', (...args) => {
        console.log('close', ...args, ws);

        emit({ type: 'healthz', payload: false });
      });

      ws.on('error', (...args) => {
        CommonLog.lError('error', ...args);

        emit({ type: 'healthz', payload: false });
        // console.log('websocket error. reconnect after 1s');
      });
      ws.on('ping', () => {
        // CommonLog.Info('receive ping msg');
      });
      ws.on('pong', () => {
        // CommonLog.Info('receive pong msg');
      });
      ws.on('message', data => {
        emit({ type: 'data', payload: data });
      });
      ws.on('websocket-status', (msg) => {
        CommonLog.Info(msg);
        // if(/Disconnected/.test(msg)){
        //   emit({ type: 'healthz', payload: false });
        // }
      })
    } else {
      emit({ type: 'healthz', payload: false });
      CommonLog.lError('ws doesn\'t exist', { at: 'createRushChannel' });
    }
    return () => {
    };
  });
}

const rushChannelHandlers = {
  healthz: rushHealthz,
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
