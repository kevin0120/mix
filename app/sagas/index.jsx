// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { cardAuthFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchOperation, watchResults } from './operation';
import { watchIO } from './io';
import { shutDownDiagWorkFlow } from './shutDownDiag';
import { toolFunctions } from './tools';
import { loginFlow, logoutFlow } from './auth';
import { healthzCheckFlow } from './healthzCheck';
import { watchAiis } from './aiis';
import { watchRush } from './rush';

export default function* rootSaga() {
  yield all([
    // card auth
    cardAuthFlow(),
    watchScanner(),
    watchAiis(),
    watchOperation(),
    watchResults(),
    watchIO(),
    watchRush(),
    shutDownDiagWorkFlow(),
    toolFunctions(),
    // auth
    loginFlow(),
    logoutFlow(),
    // healthz
    healthzCheckFlow(),
    sysInitFlow()
  ]);
}
