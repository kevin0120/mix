// @flows

import { take, put, call } from 'redux-saga/effects'

import {
  CARD_AUTH,
} from '../actions/actionTypes';

export function* getCardReaderInsert(action) {
  console.log(action);
  yield put({ type: 'INCREMENT' }); // dispatch a new action to reducer
}

export function* authFlow() {
  while (true){
    const {t } = yield take(CARD_AUTH.READER.INSERTED);
    yield call(getCardReaderInsert, t);
  }
}
