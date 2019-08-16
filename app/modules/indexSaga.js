// @flow

import { take, fork, all, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';

import sysInitFlow from './systemInit/saga';
import user from './user/saga';
import { watchAiis } from './aiis/saga';
// import { watchSettingPreSave } from './setting/saga';
import { watchRushEvent } from './rush/saga';
import watchRFIDEvent from './external/device/rfid/saga';
import watchNotification from './notification/saga';
import watchOperationViewer from './operationViewer/saga';
import watchNetwork from './network/saga';
import watchBattery from './battery/saga';
import watchPower from './power/saga';
import andon from './andon/saga';
import order from './order/saga';
import dialog from './dialog/saga';
import { CommonLog } from '../common/utils';
import healthz from './healthz/saga';
export default function* rootSaga(): Saga<void> {
  try {
    const state = yield select();
    const { andonEnable } = state.setting.systemSettings;
    yield all([
      // 硬件设备
      watchRFIDEvent(),

      watchNotification(),
      watchAiis(),
      watchRushEvent(),

      user(),// auth
      // watchSettingPreSave(),
      sysInitFlow(),
      watchOperationViewer(),
      watchNetwork(),
      watchBattery(),
      watchPower(),
      andonEnable ? andon() : null,
      order(),
      dialog(),
      healthz()
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}
