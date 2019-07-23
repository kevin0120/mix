// @flow

import { put, takeLatest, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';

import { SCANNER } from './action';
import { operationTrigger, operationTriggerBlock } from '../operation/action';
import { isCarID } from '../../common/utils';

const lodash = require('lodash');

function* scannerHandler(action) {
  try {
    const { data, source } = action;
    const state = yield select();
    if (state.setting.operationSettings.opMode === 'op'
      // && state.workMode.workMode === 'manual'
    ) {
      yield put(operationTrigger(data, null, null, source));
      return;
    }
    if (!lodash.isNil(data) && !lodash.isEmpty(data)) {
      if (isCarID(data)) {
        yield put(operationTrigger(data, null, null, source));
      } else {
        yield put(operationTrigger(null, data, null, source));
      }
      yield put(operationTriggerBlock(false));
    }
  } catch (e) {
    console.error(e);
  }
}

export default function* watchScanner(): Saga<void> {
  yield takeLatest(SCANNER.READ_NEW_DATA, scannerHandler);
}
