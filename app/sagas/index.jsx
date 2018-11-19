// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { authFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchResults } from './operation';
import { shutDownDiagWorkFlow } from './shutDownDiag';
import { toolFunctions } from './tools';

export default function* rootSaga() {
  yield all([
    sysInitFlow(),
    authFlow(),
    watchScanner(),
    watchResults(),
    shutDownDiagWorkFlow(),
    toolFunctions()
  ]);
}
