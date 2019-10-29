// @flow
import { put, all } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import NotifierActions from '../Notifier/action';
import healthzActions from '../healthz/action';
import type { tRushConnectionCallback } from './type';

const onConnectActions: Array<tRushConnectionCallback> = [
  () => NotifierActions.enqueueSnackbar('Info', 'rush 已连接')
];
const onDisconnectActions: Array<tRushConnectionCallback> = [
  () => NotifierActions.enqueueSnackbar('Error', 'rush 已断开')
];
const onChangeActions = [];

export const bindRushAction = {
  onConnect: bindOnConnectAction,
  onDisconnect: bindOnDisconnectAction,
  onChange: bindOnChangeAction
};

let rushHealthz = false;

function* bindOnConnectAction(actionToBind: tRushConnectionCallback): Saga<void> {
  try {
    onConnectActions.push(actionToBind);
    if (rushHealthz) {
      yield put(actionToBind());
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindOnDisconnectAction(actionToBind: tRushConnectionCallback): Saga<void> {
  try {
    onDisconnectActions.push(actionToBind);
    if (!rushHealthz) {
      yield put(actionToBind());
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindOnChangeAction(actionToBind: tRushConnectionCallback): Saga<void> {
  try {
    onChangeActions.push(actionToBind);
    yield put(actionToBind());
  } catch (e) {
    CommonLog.lError(e);
  }
}

// TODO make rush extends CommonExternalEntity

export default function* (payload: boolean): Saga<void> {
  try {
    yield put(healthzActions.data({ rush: payload }));
    if (rushHealthz !== payload) {
      yield all(onChangeActions.map(action => put(action(payload))));
    }
    yield all(
      (payload ? onConnectActions : onDisconnectActions).map(action =>
        put(action(payload))
      )
    );
    rushHealthz = payload;
  } catch (e) {
    CommonLog.lError(e, { at: 'rush handleHealthz' });
  }
}
