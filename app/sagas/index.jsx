// @flow


import { all } from 'redux-saga/effects'


import { watchScanner } from './scanner'
import { authFlow } from './cardAuth'

export default function * rootSaga () {
  yield all([
    authFlow(),
    watchScanner()
  ]);
}
