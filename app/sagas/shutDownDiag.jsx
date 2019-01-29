/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { select, put, call, take } from 'redux-saga/effects';

import { openShutDownDiagwMsg, closeShutDownDiag} from '../actions/shutDownDiag';

import {SHUTDOWN_DIAG } from '../actions/actionTypes'

import {
  switch2Ready,
  switch2Doing,
  operationVerified
} from '../actions/operation';
import { Info } from '../logger';
// import { toolDisable } from '../actions/tools';

import { ak2 } from './operation';
import { setNewNotification } from "../actions/notification";

const { ipcRenderer } = require('electron');

// const getOperations = state => state.operations;

function* openDiag(dType, data) {
  const state = yield select();
  const { enable } = state.setting.operationSettings.byPass;
  if (dType === 'bypass' && !enable) {
    // bypass 但是其功能未使能
    return;
  }
  const msg = data !== null ? data.reasons.join(',') : '';
  yield put(openShutDownDiagwMsg(dType,msg));
}

function* closeDiag(dType) {
  try {
    if (dType === 'verify') {
      yield put(switch2Ready());
    }

    yield put(closeShutDownDiag());
  }catch (e) {
    console.error(`closeDiag error: ${e.message}`)
  }

}

function* confirmDiag(dType, data) {

  try {

    const state = yield select();

    switch (dType) {
      case 'shutdown': {
        ipcRenderer.send('asynchronous-message', 'shutdown');
        break;
      }
      case 'bypass': {
        const {operations:op } = state;
        const { carID } = op;
        const { enableAk2=true } = state.setting.operationSettings;
        if (enableAk2) {
          yield call(ak2);
        }
        Info(`车辆已放行 车辆ID:${carID}`);
        yield put(switch2Ready());
        break;
      }
      case 'verify': {
        // 冲突确认，继续作业
        const {enableConflictOP=false} = state.setting.systemSettings;
        if(enableConflictOP){
          yield put(operationVerified(data));
        }else {
          yield put(setNewNotification('warning', `设定为不允许重复拧紧同一张工单 VIN: ${data.vin}`));
          return; // 直接返回, 不关闭模式对话框
        }
        break;
      }
      default: {
        yield put(switch2Doing());
      }
    }

    yield put(closeShutDownDiag());
  } catch (e) {
    console.error(`confirmDiag error: ${e.message}`)
  }
}

export function* shutDownDiagWorkFlow() {
  try {
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
  }catch (e) {
    console.error(`shutDownDiagWorkFlow error: ${e.message}`)
  }

}
