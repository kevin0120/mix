// @flow

import { take, call, all, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';

import sysInitFlow from './systemInit/saga';
import user from './user/saga';
import { watchAiis } from './aiis/saga';
// import { watchSettingPreSave } from './setting/saga';
import { watchRushEvent } from './rush/saga';
import watchRFIDEvent from './external/device/rfid/saga';
import watchOperationViewer from './operationViewer/saga';
import watchNetwork from './network/saga';
import watchBattery from './battery/saga';
import watchPower from './power/saga';
import andon from './andon/saga';
import order from './order/saga';
import dialog from './dialog/saga';
import { CommonLog } from '../common/utils';
import healthz from './healthz/saga';
import modelViewer from './modelViewer/saga';
import notifier from './Notifier/saga';

export default function* rootSaga(): Saga<void> {
  try {
    const state = yield select();
    const { andonEnable } = state.setting.systemSettings;
    yield all([
      // 硬件设备
      watchRFIDEvent,
      watchAiis,
      watchRushEvent,

      user, // auth
      // watchSettingPreSave(),
      sysInitFlow,
      watchOperationViewer,
      // watchNetwork(),
      // watchBattery(),
      watchPower,
      andonEnable ? andon : null,
      order,
      dialog,
      healthz,
      modelViewer,
      notifier
    ].filter(e => !!e).map(e => e && call(e)));
  } catch (e) {
    CommonLog.lError(e);
  }
}
