import { take, put, select, join, fork, cancel } from 'redux-saga/effects';
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
      const step = yield select(s => processingStep(s.order));
      const action = yield take(SCANNER_STEP.GET_VALUE);
      yield put(orderActions.stepData((data) => ({
        ...data,
        result: {
          [step.payload.label]: action.value
        }
      })));
      yield put(orderActions.stepStatus(STEP_STATUS.LEAVING));
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.LEAVING](ORDER, orderActions) {
    try {
      while (true) {
        yield take(SCANNER_STEP.SUBMIT);
        const result = yield select(s => stepData(processingStep(s.order))?.result);
        const label = yield select(s => stepPayload(processingStep(s.order))?.label);
        if (result.hasOwnProperty(label)) {
          yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
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

