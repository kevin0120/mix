import { put, take } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import { VIDEO_STEP } from './action';

const videoStepMixin = (BaseStep) => class ClsVideoStep extends BaseStep {
  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        yield take(VIDEO_STEP.SUBMIT);
        yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
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
    }
  };
};
export default videoStepMixin;
