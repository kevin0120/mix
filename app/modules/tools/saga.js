// @flow

import { call, select, put } from 'redux-saga/effects';

import { TOOLS } from './action';

import { toolEnable } from '../../api/operation';
import { watchWorkers } from '../util';
import { setNewNotification } from '../notification/action';
import ClsScrewTool, { defaultScrewToolDispatcher } from './model';

type actionType = {
  +type: string,
  +enable: boolean
};

type controllerType = {
  +connection: string,
  +serial_no: string
};

export const staticScrewTool = new ClsScrewTool('G1', "0001");
staticScrewTool.dispatcher = defaultScrewToolDispatcher;

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
    yield call(toolEnable, mUrl, targetResult.controller_sn, targetResult.gun_sn, action.enable, action.reason);
  } catch (e) {
    console.error(`staticToolEnable error:`, e);
    const { data } = e.response || { data: '' };
    if (/tool not found/.test(data)) {
      const state = yield select();
      const { results } = state.operations;
      const targetResult = results[0];

      yield put(setNewNotification('Warn', `工具序列号不匹配：${targetResult.gun_sn}`));
    }
  }

}

function* onToolStatusChange(action) {
  try {
    const { toolSN, status, reason } = action;
    yield put(
      setNewNotification(
        'info',
        `拧紧枪状态更新（${toolSN}）： ${status}${reason ? `, ${reason}` : ''}`
      )
    );
  } catch (e) {
    console.error('onToolStatusChange:', e);
  }
}

function* onToolResult(action) {
  try {
    const { results } = action;
    staticScrewTool.doDispatch(results);
  } catch (e) {
    console.error('onToolResult:', e);
  }
}

const workers = {
  [TOOLS.ENABLE]: staticToolEnable,
  [TOOLS.STATUS_CHANGE]: onToolStatusChange,
  [TOOLS.NEW_RESULTS]: onToolResult
};

export const toolFunctions = watchWorkers(workers);

// export function* toolFunctions() {
//   yield takeEvery(TOOLS.ENABLE, staticToolEnable);
// }
