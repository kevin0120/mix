// @flow

import { take, call } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';
import { triggerOperation } from './operation';
import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';

export function* watchScanner() {
  while (true) {
    const { data } = yield take(SCANNER.READ_NEW_DATA);
    if (isCarID(data)) {
      yield call(triggerOperation, data, null, null, OPERATION_SOURCE.SCANNER);
    } else {
      yield call(triggerOperation, null, data, null, OPERATION_SOURCE.SCANNER);
    }
  }
}
