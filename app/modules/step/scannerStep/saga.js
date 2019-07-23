import { take, put, select } from 'redux-saga/effects';
import { SCANNER_STEP } from './action';
import { processingStep } from '../../order/selector';

export default function* root(ORDER, orderActions) {
  try {
    while (true) {
      const step=yield select(s=>processingStep(s.order));
      const action = yield take(SCANNER_STEP.GET_VALUE);
      yield put(orderActions.stepData((data)=>({
        ...data,
        result:{
          [step.payload.label]:action.value
        }
      })))
    }
  } catch (e) {
    console.error(e);
  }
}
