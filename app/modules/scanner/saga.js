// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import {SCANNER} from './action';
import Scanner from './model';

let scanner = new Scanner('Scanner');

function* scannerHandler(action): Saga<void> {
  try {
    const { data } = action;
    if (scanner.validate(data)){
      yield put(scanner.dispatch())
    }else {
      // do nothing
    }
  } catch (e) {
    console.error(e);
  }
}

export default function* watchScanner(): Saga<void> {
  yield takeLatest(SCANNER.READ_NEW_DATA, scannerHandler);
}
