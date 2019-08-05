// @flow

import { take, fork, all, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';

import watchScannerEvent, { scanner } from './scanner/saga';
import watchReader, { reader } from './reader/saga';
import watchController from './controller/saga';
import { cardAuthFlow } from './cardAuth/saga';
import sysInitFlow from './systemInit/saga';
import { operationFlow, watchResults } from './operation/saga';
import { watchIOEvent } from './io/saga';
import { toolFunctions } from './tools/saga';
import user from './user/saga';
import { watchAiis } from './aiis/saga';
import { watchSettingPreSave } from './setting/saga';
import { watchRushEvent } from './rush/saga';
import watchRFIDEvent from './rfid/saga';
import watchNotification from './notification/saga';
import watchOperationViewer from './operationViewer/saga';
import logoFlow from './logo/saga';
import watchNetwork from './network/saga';
import watchBattery from './battery/saga';
import watchPower from './power/saga';
import andon from './andon/saga';
import order from './order/saga';
import dialog from './dialog/saga';
import { CommonLog } from '../common/utils';

export default function* rootSaga(): Saga<void> {
  try {
    const state = yield select();
    const { andonEnable } = state.setting.systemSettings;
    yield all([
      // card auth
      // cardAuthFlow(),
      // 硬件设备
      watchScannerEvent(),
      watchReader(),
      watchController(),
      watchIOEvent(), // 监听IO数据
      watchRFIDEvent(),

      watchNotification(),
      watchAiis(),
      operationFlow(),
      watchResults(), // 监听结果
      watchRushEvent(),
      toolFunctions(),
      user(),// auth
      // healthz
      // healthzCheckFlow(),
      // watchSettingPreSave(),
      sysInitFlow(),
      watchOperationViewer(),
      // logoFlow(),
      watchNetwork(),
      watchBattery(),
      watchPower(),
      andonEnable ? andon() : null,
      order(),
      dialog()
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}
