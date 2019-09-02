import { put, take } from 'redux-saga/effects';
import Step from '../Step';
import STEP_STATUS from '../model';
import { INPUT_STEP } from './action';



export default class InputStep extends Step {
  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        yield put(orderActions.stepStatus(this,STEP_STATUS.DOING));
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        while (true) {
          const { payload } = yield take(INPUT_STEP.SUBMIT);
          if (payload) {
            if (payload === 'fail') {
              yield put(orderActions.stepStatus(this,STEP_STATUS.FAIL));
            } else {
              yield put(orderActions.stepStatus(this,STEP_STATUS.FINISHED));
            }
          }
        }
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.FAIL](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        console.error(e);
      }
    }
  };
}
