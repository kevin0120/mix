import { put, take } from 'redux-saga/effects';
import { INPUT_STEP } from './action';

export default function* root(ORDER, orderActions) {
  try {
    console.log('input running');
    while (true) {
      const { payload } = yield take(INPUT_STEP.SUBMIT);
      if (payload) {
        if (payload === 'fail') {
          yield put(orderActions.failStep());
        } else {
          yield put(orderActions.finishStep());
        }
        yield put(orderActions.pushStep());
      }
    }
  } catch (e) {
    console.error(e);
  }
}
