// @flow

import { put, takeLatest } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';
import { operationTrigger, operationTriggerBlock } from '../actions/operation';
// import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';

const lodash = require('lodash');

function* scannerHandler(action) {
  try {
    const { data, source } = action;
    if (!lodash.isNil(data) && data !== '') {
      if (isCarID(data)) {
        yield put(operationTrigger( data, null, null, source));
      } else {
        yield put(operationTrigger( null, data, null, source));
      }
      yield put(operationTriggerBlock(false));

    }
  } catch (e) {
    console.log(e.message);
  }
}

export function* watchScanner() {
  yield takeLatest(SCANNER.READ_NEW_DATA, scannerHandler);
}
