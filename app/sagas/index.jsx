// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { authFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchOperation, watchResults } from './operation';
import { watchIO } from './io';
import { watchHealth} from './health';
import { shutDownDiagWorkFlow } from './shutDownDiag';
import { toolFunctions } from './tools';

export default function* rootSaga() {
  yield all([
    authFlow(),
    watchScanner(),
    watchResults(),
    watchIO(),
    watchHealth(),
    shutDownDiagWorkFlow(),
    toolFunctions(),
    watchOperation(),
    sysInitFlow(),
  ]);
}
