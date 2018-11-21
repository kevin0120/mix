// @flow

import { take, call } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';
import { triggerOperation } from './operation';
import { OPERATION_SOURCE } from '../reducers/operations';

export function* watchScanner() {
  while (true) {
    const { data } = yield take(SCANNER.READ_NEW_DATA);
    yield call(triggerOperation, data, OPERATION_SOURCE.SCANNER);
  }
}
