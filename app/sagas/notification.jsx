/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { takeEvery, select } from 'redux-saga/effects';

import { NOTIFY } from '../actions/actionTypes';
import { Info, Warn, Error, Maintenance } from '../logger';

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
