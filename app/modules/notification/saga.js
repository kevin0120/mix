// @flow

import { takeEvery } from 'redux-saga/effects';

import { NOTIFY } from './action';
import { Info, Warn, Error, Maintenance } from '../../logger';

function* notificationAlways(action) {
  const { variant, message, meta } = action;
  switch (variant) {
    case 'info':
      Info(message,meta);
      break;
    case 'warning':
      Warn(message,meta);
      break;
    case 'maintenance':
      Maintenance(message,meta);
      break;
    case 'error':
      Error(message,meta);
      break;
    default:
      break;
  }
}

export function* watchNotification() {
  yield takeEvery(NOTIFY.NEW_NOTIFICATION, notificationAlways);
}
