// @flow

import type { Saga } from 'redux-saga';
import { takeEvery, call } from 'redux-saga/effects';

import { NOTIFY } from './action';
import { Info, Warn, lError, Maintenance } from '../../logger';
import type { tNotifyVariant } from './action';
import { CommonLog } from '../../common/utils';

type tNotifyFuncMap ={
  [type: tNotifyVariant]: (string | Error, string) => void
};

const notifyFuncMap: tNotifyFuncMap = {
  'Info': Info,
  'Warn': Warn,
  'Maintenance': Maintenance,
  'Error': lError
};

function* notificationAlways(action): Saga<void> {
  try {
    const { variant = 'Info', message, meta } = action;
    const t = (variant: tNotifyVariant);
    const method = notifyFuncMap?.[t];
    if (!method) return;
    yield call(method, message, meta);
  } catch (e) {
    CommonLog.lError(e);
  }

}

export default function* watchNotification(): Saga<void> {
  try {
    yield takeEvery(NOTIFY.NEW_NOTIFICATION, notificationAlways);
  } catch (e) {
    CommonLog.lError(e);
  }
}
