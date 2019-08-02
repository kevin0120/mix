/* eslint-disable no-bitwise,no-case-declarations */
// @flow
import { put, call, fork } from 'redux-saga/effects';
import { channel } from 'redux-saga';
import { loginRequest, logoutRequest } from '../user/action';
import { CARD_AUTH } from './action';

export function* cardAuthFlow() {
  try {
    pcscd = yield call(pcscLite);
    yield call(pcscdListener);
  } catch (e) {
    console.error(e);
  } finally {
    console.log('cardAuthFlow finished');
  }
}
