/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { delay } from 'redux-saga';
import { call, take, select, all, put, fork, cancel } from 'redux-saga/effects';
import { HEALTHZ_CHECK } from '../actions/actionTypes';
import { setHealthzCheck } from '../actions/healthCheck';
import { masterPCHealthCheck, controllerHealthCheck } from './api/healthzCheck';
import { setNewNotification } from '../actions/notification';

const lodash = require('lodash');

const getHealthz = state => state.healthCheckResults;

function* healthzCheckTask(url, controllers) {
  while (true) {
    try {
      const [mHealthz, cHealthzs] = yield all([
        call(masterPCHealthCheck, url),
        call(controllerHealthCheck, url, controllers)
      ]);
      const healthzStatus = yield select(getHealthz); // 获取整个healthz
      const s = mHealthz.status === 204;
      if (
        !lodash.isEqual(
          healthzStatus.masterpc.isHealth,
          s
        )
      ) {
        // 如果不相等 更新
        yield put(setHealthzCheck('masterpc', s));
        yield put(setNewNotification('info', `masterPC连接状态更新: ${s}`));
      }
      // 控制器healthz check
      const statusCode = cHealthzs.status;
      let controllerHealthzStatus = false;
      if (statusCode === 200) {
        controllerHealthzStatus = cHealthzs.data[0].status === 'online';
      }
      if (
        !lodash.isEqual(
          healthzStatus.controller.isHealth,
          controllerHealthzStatus
        )
      ) {
        yield put(setHealthzCheck('controller', controllerHealthzStatus));
        yield put(setNewNotification('info', `控制器连接状态更新: ${controllerHealthzStatus}`));
      }
    } catch (e) {
      // yield put(setNewNotification('error', e.toString()));
    }

    yield call(delay, 3000); // 延时 3s
  }
}

function* getConnectionInfo() {
  const state = yield select();
  const u = state.connections.masterpc;
  const c = state.connections.controllers.map(v => v.serial_no);
  return { u, c };
}

function* startHealthzCheck(url, controllers) {
  let U = url;
  let C = controllers;
  if (lodash.isNil(url) || lodash.isNil(controllers)) {
    const { u, c } = yield call(getConnectionInfo);
    U = u;
    C = c;
  }
  yield call(healthzCheckTask, U, C);
}

function* restartHealthzCheck() {
  const { u, c } = yield call(getConnectionInfo);
  yield put({ type: HEALTHZ_CHECK.START, u, c }); // 启动healthzCheck
}

export function* healthzCheckFlow() {
  while (true) {
    const { url, controllers } = yield take(HEALTHZ_CHECK.START);
    const task = yield fork(startHealthzCheck, url, controllers);
    yield take(HEALTHZ_CHECK.RESTART);
    yield cancel(task);
    yield call(restartHealthzCheck);
  }
}
