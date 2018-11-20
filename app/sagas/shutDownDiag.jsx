/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { select, put, call, take } from 'redux-saga/effects';

import { SHUTDOWN_DIAG } from '../actions/actionTypes';

import { forceSwitch2Ready, switch2Doing } from '../actions/operation';
import { Info } from '../logger';
import { toolDisable } from '../actions/tools';

const { ipcRenderer } = require('electron');

const getOperations = state => state.operations;

export function* openDiag(dType, data) {
  const state = yield select();
  const { enable } = state.operationSettings.byPass;
  if (dType === 'bypass' && !enable) {
    // bypass 但是其功能未使能
    return;
  }
  const msg = data !== null ? data.reasons.join(',') : '';
  yield put({ type: SHUTDOWN_DIAG.OPEN_WITH_MSG, dType, msg });
}

export function* closeDiag(dType) {
  if (dType === 'verify') {
    yield put(forceSwitch2Ready());
  }

  yield put({ type: SHUTDOWN_DIAG.CLOSE });
}

export function* confirmDiag(dType, data) {
  switch (dType) {
    case 'shutdown': {
      ipcRenderer.send('asynchronous-message', 'shutdown');
      break;
    }
    case 'bypass': {
      const op = yield select(getOperations);
      const { carID } = op;
      Info(`车辆已放行 车辆ID:${carID}`);
      yield put(toolDisable());
      yield put(forceSwitch2Ready());
      break;
    }
    case 'verify': {
      break;
    }
    default: {
      yield put(switch2Doing());
    }
  }

  yield put({ type: SHUTDOWN_DIAG.CLOSE });
}

export function* shutDownDiagWorkFlow() {
  while (true) {
    const { dType, data } = yield take(SHUTDOWN_DIAG.OPEN);
    yield call(openDiag, dType, data);
    const action = yield take([
      SHUTDOWN_DIAG.CLOSE_START,
      SHUTDOWN_DIAG.CONFIRM
    ]);
    switch (action.type) {
      case SHUTDOWN_DIAG.CLOSE_START:
        yield call(closeDiag, dType);
        break;
      case SHUTDOWN_DIAG.CONFIRM:
        yield call(confirmDiag, dType, data);
        break;
      default:
        break;
    }
  }
}
