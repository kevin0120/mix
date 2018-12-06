// @flow

import { take, call } from 'redux-saga/effects';

import { SCANNER } from '../actions/actionTypes';
import { triggerOperation } from './operation';
// import { OPERATION_SOURCE } from '../reducers/operations';
import { isCarID } from '../common/utils';

const lodash = require('lodash');

export function* watchScanner() {
  while (true) {
    const { data, source } = yield take(SCANNER.READ_NEW_DATA);
    if ( !lodash.isNil(data) && data !== ''){
      if (isCarID(data)) {
        yield call(triggerOperation, data, null, null, source);
      } else {
        yield call(triggerOperation, null, data, null, source);
      }
    }
  }
}
