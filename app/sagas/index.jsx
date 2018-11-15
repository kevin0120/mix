// @flow

import { all } from 'redux-saga/effects';

import { watchScanner } from './scanner';
import { authFlow } from './cardAuth';
import { sysInitFlow } from './systemInit';

export default function* rootSaga() {
  yield all([sysInitFlow(), authFlow(), watchScanner()]);
}
