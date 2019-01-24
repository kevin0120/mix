// @flow

import { call, takeLatest } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';
import { triggerOperation } from './operation';
// import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';

const lodash = require('lodash');

function* scannerHandler(action) {
  try {
    const { data, source } = action;
    if (!lodash.isNil(data) && data !== '') {
      if (isCarID(data)) {
        yield call(triggerOperation, data, null, null, source);
      } else {
        yield call(triggerOperation, null, data, null, source);
      }
    }
  } catch (e) {
    console.log(e.message);
  }
}

export function* watchScanner() {
  yield takeLatest(SCANNER.READ_NEW_DATA, scannerHandler);
}
