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
  try {
    const state = yield select();
    const mUrl = state.connections.masterpc;
    // const controller: controllerType = state.connections.controllers[0];
    // if (controller === undefined) {
    //   return;
    // }

    // const toolSN = state.setting.systemSettings.defaultToolSN || '';

    const { results } = state.operations;
    const targetResult = results[0];

    yield call(toolEnable, mUrl, targetResult.controller_sn, targetResult.gun_sn, action.enable);
  }catch (e) {
    console.error(`staticToolEnable error: ${e.message}`);
  }

}

export function* toolFunctions() {
  yield takeEvery(TOOLS.ENABLE, staticToolEnable);
}
