// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { cardAuthFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchOperation, watchResults } from './operation';
import { watchIO } from './io';
import { toolFunctions } from './tools';
import { loginFlow, logoutFlow } from './auth';
// import { healthzCheckFlow } from './healthzCheck';
import { watchAiis } from './aiis';
import { watchSettingPreSave } from './setting';
import { watchRush } from './rush';
import { watchRfid } from './rfid';
import { watchNotification } from './notification';
import watchOperationViewer from './operationViewer';
import logo from './logo';
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
      watchOperation(),
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
      logo(),
      watchNetwork(),
      watchBattery(),
      watchPower()
    ]);
  } catch (e) {
    console.log('rootSaga:', e);
  }
}
