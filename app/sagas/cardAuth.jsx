// @flow

import { takeLatest, put, call } from 'redux-saga/effects';
import { CARD_AUTH, USER, USER_CONFIGS } from '../actions/actionTypes';
import { doUserAuth, userLogOut } from '../actions/userAuth';

export function* cardAuthFlow() {
  yield takeLatest(CARD_AUTH.CARD.READ.SUCCESS, cardReadSuccess);
  yield takeLatest(CARD_AUTH.READER.REMOVED, cardRemoved);
  yield takeLatest(CARD_AUTH.CARD.REMOVED, cardRemoved);
}

export function* cardReadSuccess(action) {
  const { data } = action;
  yield put(doUserAuth(data));
}

export function* cardRemoved() {
  yield put(userLogOut('112233')); //默认uuid112233 登出
}
