import { take, put } from 'redux-saga/effects';
import { SCANNER_STEP } from './action';

export default function* root(ORDER, orderActions) {
  try {
    while (true) {
      const action = yield take(SCANNER_STEP.GET_VALUE);
      yield put(orderActions.finishStep());
    }
  } catch (e) {
    console.error(e);
  }
}
