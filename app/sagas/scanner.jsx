// @flow

import { select, put, takeEvery } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';

export function* incrementAsync(action) {
  console.log(action);
  const state = yield select();
  console.log(state);
  yield put({ type: 'INCREMENT' });
}

export function* watchScanner() {
  yield takeEvery(SCANNER.READ_NEW_DATA, incrementAsync);
}
