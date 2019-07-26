import { take, put, select } from 'redux-saga/effects';
import { SCANNER_STEP } from './action';
import { processingStep, stepData, stepPayload } from '../../order/selector';
import { STEP_STATUS } from '../model';

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      while (true) {
        const action = yield take([SCANNER_STEP.GET_VALUE, SCANNER_STEP.SUBMIT]);
        const result = yield select(s => stepData(processingStep(s.order))?.result);
        const label = yield select(s => stepPayload(processingStep(s.order))?.label);
        switch (action.type) {
          case(SCANNER_STEP.GET_VALUE):
            yield put(orderActions.stepData((data) => ({
              ...data,
              result: {
                [label]: action.value
              }
            })));
            break;
          case(SCANNER_STEP.SUBMIT):
            if (Object.hasOwnProperty.call(result || {}, label)) {
              yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
            }
            break;
          default:
            break;
        }
      }
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      console.error(e);
    }
  }
};

