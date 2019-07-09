import { put, take } from 'redux-saga/effects';
import {INPUT_STEP}from './action';

export default function* root(ORDER, orderActions) {
  try {
    console.log('input running');
    while(true){
      const { payload } = yield take(INPUT_STEP.SUBMIT);
      if (payload) {
        yield put(orderActions.pushStep());
      }
    }
  } catch (e) {
    console.error(e);
  }
}
