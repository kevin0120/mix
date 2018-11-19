// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { authFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';
import { watchResults } from './operation';
import { watchIO } from './io';
import { watchHealth} from './health';

export default function * rootSaga () {
  yield all([
    authFlow(),
    watchScanner(),
    watchResults(),
    watchIO(),
    watchHealth(),
    sysInitFlow(),
  ]);
}
