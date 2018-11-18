// @flow

import { select, put, take, call } from 'redux-saga/effects';

import { SCANNER, OPERATION } from '../actions/actionTypes';
import { getOperation } from './operation';

export function* watchScanner() {
  while (true){
    const { data } = yield take(SCANNER.READ_NEW_DATA);
    yield put({type: OPERATION.TRIGGER.NEW_DATA, data: data});

    const state = yield select();

    const triggers = state.setting.operationSettings.flowTriggers;

    let triggerFlagNum = 0;
    for (const i in triggers) {
      if (state.operations[triggers[i]] !== '') {
        triggerFlagNum += 1;
      }
    }

    if (triggerFlagNum === triggers.length) {
      yield call(getOperation);
    }
  }
}
