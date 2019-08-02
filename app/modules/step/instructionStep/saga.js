import { put, take } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import { INSTRUCTION_STEP } from './action';

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
      yield take(INSTRUCTION_STEP.SUBMIT);
      yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
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
