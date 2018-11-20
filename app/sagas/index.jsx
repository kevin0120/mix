// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
// import { cardAuthFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchResults } from './operation';
import { shutDownDiagWorkFlow } from './shutDownDiag';
import { toolFunctions } from './tools';
import { loginFlow, logoutFlow } from './auth';
import { healthzCheckFlow } from './healthzCheck';

export default function* rootSaga() {
  yield all([
    sysInitFlow(),
    // card auth
    // cardAuthFlow(),
    watchScanner(),
    watchResults(),
    shutDownDiagWorkFlow(),
    toolFunctions(),
    // auth
    loginFlow(),
    logoutFlow(),
    // healthz
    healthzCheckFlow()
  ]);
}
