// @flow


import { all } from 'redux-saga/effects'


import { watchScanner } from './scanner'
import { authFlow } from './cardAuth'
import { sysInitFlow } from './systemInit'
import { watchResults } from './operation';

export default function * rootSaga () {
  yield all([
    sysInitFlow(),
    authFlow(),
    watchScanner(),
    watchResults()
  ]);
}
