/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { call, takeEvery, select } from 'redux-saga/effects';

import { TOOLS } from '../actions/actionTypes';

import { toolEnable } from './api/operation';

type actionType = {
  +type: string,
  +enable: boolean
};

type controllerType = {
  +connection: string,
  +serial_no: string
};

function* staticToolEnable(action: actionType) {
  const state = yield select();
  const mUrl = state.connections.masterpc;
  const controller: controllerType = state.connections.controllers[0];
  if (controller === undefined) {
    return;
  }

  yield call(toolEnable, mUrl, controller.serial_no, action.enable);
}

export function* toolFunctions() {
  yield takeEvery(TOOLS.ENABLE, staticToolEnable);
}
