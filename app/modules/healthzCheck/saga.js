// @flow

import { call, take, select, all, put, fork, cancel, delay } from 'redux-saga/effects';
import { HEALTHZ_CHECK, setHealthzCheck } from './action';
import { masterPCHealthCheck, controllerHealthCheck } from '../../api/healthzCheck';
import { setNewNotification } from '../notification/action';

const lodash = require('lodash');

const getHealthz = state => state.healthCheckResults;

let task = null;

function* healthzCheckTask(url, controllers) {
  try {
    while (true) {
      try {
        const [mHealthz, cHealthzs] = yield all([
          call(masterPCHealthCheck, url),
          call(controllerHealthCheck, url, controllers)
        ]);
        const healthzStatus = yield select(getHealthz); // 获取整个healthz
        const s = mHealthz.status === 204;
        if (!lodash.isEqual(healthzStatus.masterpc.isHealth, s)) {
          // 如果不相等 更新
          yield put(setHealthzCheck('masterpc', s));
          yield put(setNewNotification('Info', `masterPC连接状态更新: ${s}`));
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
          yield put(
            setNewNotification(
              'info',
              `控制器连接状态更新: ${controllerHealthzStatus}`
            )
          );
        }
      } catch (e) {
        console.error(e);
        // yield put(setNewNotification('error', e.toString()));
      }

      yield delay(3000); // 延时 3s
    }
  } catch (e) {
    console.error(e);
  }

}

function* getConnectionInfo() {
  try {
    const state = yield select();
    const u = state.connections.masterpc;
    const c = state.connections.controllers.map(v => v.serial_no);
    return { u, c };
  } catch (e) {
    console.error(e);
  }
}

function* startHealthzCheck(url, controllers) {
  try {
    if (!lodash.isNil(task)) {
      yield cancel(task);
    }
    let U = url;
    let C = controllers;
    if (lodash.isNil(url) || lodash.isNil(controllers)) {
      const { u, c } = yield call(getConnectionInfo);
      U = u;
      C = c;
    }
    yield call(healthzCheckTask, U, C);
  } catch (e) {
    console.error(e);
  }
}

//
// function* restartHealthzCheck() {
//   const { u, c } = yield call(getConnectionInfo);
//   yield put({ type: HEALTHZ_CHECK.START, u, c }); // 启动healthzCheck
// }

export default function* healthzCheckFlow() {
  try {
    while (true) {
      const { url, controllers } = yield take(HEALTHZ_CHECK.START);
      task = yield fork(startHealthzCheck, url, controllers);
    }
  } catch (e) {
    console.error(e);
  }
}
