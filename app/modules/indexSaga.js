// @flow

import { take, fork, all, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';

import watchScanner from './scanner/saga';
import { cardAuthFlow } from './cardAuth/saga';
import sysInitFlow from './systemInit/saga';
import { operationFlow, watchResults } from './operation/saga';
import { watchIO } from './io/saga';
import { toolFunctions } from './tools/saga';
import user from './user/saga';
import { watchAiis } from './aiis/saga';
import { watchSettingPreSave } from './setting/saga';
import { watchRush } from './rush/saga';
import { watchRfid } from './rfid/saga';
import { watchNotification } from './notification/saga';
import watchOperationViewer from './operationViewer/saga';
import logoFlow from './logo/saga';
import watchNetwork from './network/saga';
import watchBattery from './battery/saga';
import watchPower from './power/saga';
import andon from './andon/saga';
import order from './order/saga';
import dialog from './dialog/saga';
// import healthzCheckFlow from './healthzCheck';

export function watch(workers, channel): Saga<void> {
  return function* watcher() {
    try {
      while (true) {
        const action = yield take(channel || Object.keys(workers));
        if (workers[action.type].length > 1) {
          const effect = workers[action.type][0];
          const worker = workers[action.type][1];
          yield effect(worker, action);
        } else {
          yield fork(workers[action.type], action);
        }
      }
    } catch (e) {
      console.error('saga watcher error:', e, workers, channel);
    }
  };
}

export default function* rootSaga(): Saga<void> {
  try {
    const state = yield select();
    const { andonEnable } = state.setting.systemSettings;
    yield all([
      // card auth
      // cardAuthFlow(),
      watchScanner(),
      watchNotification(),
      watchAiis(),
      operationFlow(),
      watchResults(), // 监听结果
      watchIO(), // 监听IO数据
      watchRush(),
      watchRfid(),
      toolFunctions(),
      user(),// auth
      // healthz
      // healthzCheckFlow(),
      // watchSettingPreSave(),
      sysInitFlow(),
      watchOperationViewer(),
      logoFlow(),
      watchNetwork(),
      watchBattery(),
      watchPower(),
      andonEnable ? andon() : null,
      order(),
      dialog()
    ]);
  } catch (e) {
    console.error('rootSaga:', e);
  }
}
