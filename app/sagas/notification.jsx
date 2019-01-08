/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { takeEvery, put } from 'redux-saga/effects';

import { NOTIFY } from '../actions/actionTypes';
import { NewNotificationOK } from '../actions/notification';
import { Info, Warn, Error, Maintenance } from '../logger';

function* notificationAlways(action) {
  const { variant, message } = action;
  yield put(NewNotificationOK(variant, message));
  switch (variant) {
    case 'info':
      Info(message);
      break;
    case 'warning':
      Warn(message);
      break;
    case 'maintenance':
      Maintenance(message);
      break;
    case 'error':
      Error(message);
      break;
    default:
      break;
  }
}

export function* watchNotification() {
  yield takeEvery(NOTIFY.PRE_NEW_NOTIFICATION, notificationAlways);
}
