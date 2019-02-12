// @flow

import { all } from 'redux-saga/effects';

import watchScanner from './scanner';
import { cardAuthFlow } from './cardAuth';
import sysInitFlow from './systemInit';
import { operationFlow, watchResults } from './operation';
import { watchIO } from './io';
import { toolFunctions } from './tools';
import { loginFlow, logoutFlow } from './auth';
// import healthzCheckFlow from './healthzCheck';
import { watchAiis } from './aiis';
import { watchSettingPreSave } from './setting';
import { watchRush } from './rush';
import { watchRfid } from './rfid';
import { watchNotification } from './notification';
import watchOperationViewer from './operationViewer';
import logoFlow from './logo';
import watchNetwork from './network';
import watchBattery from './battery';
import watchPower from './power';

export default function* rootSaga() {
  try {
    yield all([
      // card auth
      cardAuthFlow(),
      watchScanner(),
      watchNotification(),
      watchAiis(),
      operationFlow(),
      watchResults(), // 监听结果
      watchIO(), // 监听IO数据
      watchRush(),
      watchRfid(),
      toolFunctions(),
      // auth
      loginFlow(),
      logoutFlow(),
      // healthz
      // healthzCheckFlow(),
      watchSettingPreSave(),
      sysInitFlow(),
      watchOperationViewer(),
      logoFlow(),
      watchNetwork(),
      watchBattery(),
      watchPower(),
    ]);
  } catch (e) {
    console.error('rootSaga:', e);
  }
}
