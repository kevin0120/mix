// @flow
import { call, fork, cancel, takeEvery } from 'redux-saga/effects';
import { eventChannel } from 'redux-saga';
import type { Saga, EventChannel } from 'redux-saga';
import { RUSH } from './action';
import { CommonLog } from '../../common/utils';
import handleData from './handleData';
import rushHealthz from './rushHealthz';

let task = null;

// const { getWSClient } = require('electron').remote.require('./main/webSocket');
const { ipcRenderer } = require('electron');
const { getGlobal } = require('electron').remote;

const getWSClient = getGlobal('getWSClient');

export function* watchRushEvent(): Saga<void> {
  try {
    yield takeEvery(RUSH.INIT, initRush);
  } catch (e) {
    CommonLog.lError(e, { at: 'watchRushEvent' });
  }
}

function* initRush() {
  try {
    while (!(getWSClient() && task)) {
      Object.keys(listeners).forEach(k => ipcRenderer.removeAllListeners(k));
      if (task) {
        yield cancel(task);
      }
      task = yield fork(watchRushChannel);
      ipcRenderer.send('rush');
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'initRush' });
  } finally {
    console.log('rush init finished');
  }
}

const listeners = {
  'rush-open': emit => () => {
    emit({ type: 'healthz', payload: true });
  },
  'rush-close': emit => () => {
    emit({ type: 'healthz', payload: false });
  },
  'rush-error': emit => (ev, arg) => {
    CommonLog.lError('error', arg);
    emit({ type: 'healthz', payload: false });
  },
  'rush-message': emit => (ev, data) => {
    emit({ type: 'data', payload: data });
  },
  'rush-status': emit => (ev, msg) => {
    CommonLog.Info(msg);
  }
};

function createRushChannel(): EventChannel<void> {
  return eventChannel(emit => {
    Object.keys(listeners).forEach(k => {
      ipcRenderer.on(k, listeners[k](emit));
    });
    return () => {};
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

export function* watchRushChannel(): Saga<void> {
  try {
    const chan = yield call(createRushChannel);
    yield takeEvery(chan, rushChannelTask);
  } catch (e) {
    CommonLog.lError(e, { at: 'watchRushChannel' });
  }
}
