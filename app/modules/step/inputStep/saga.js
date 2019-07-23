import { put, take } from 'redux-saga/effects';
import { INPUT_STEP } from './action';
import { STEP_STATUS } from '../../order/model';

export default function* root(ORDER, orderActions) {
  try {
    console.log('input running');
    while (true) {
      const { payload } = yield take(INPUT_STEP.SUBMIT);
      if (payload) {
        if (payload === 'fail') {
          yield put(orderActions.stepStatus(STEP_STATUS.FAIL));
        } else {
          yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
        }
      }
      yield put(orderActions.doNextStep());
    }
  } catch (e) {
    console.error(e);
  }
}
