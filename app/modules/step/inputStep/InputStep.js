// @flow
import { put, take } from 'redux-saga/effects';
import {STEP_STATUS} from '../constants';
import { INPUT_STEP } from './action';
import type { IWorkStep } from '../interface/IWorkStep';
import { CommonLog } from '../../../common/utils';
import {orderActions} from '../../order/action';

function inputStepStatusMixin(superTasks){
  return{
    ...superTasks,
    *[STEP_STATUS.DOING]() {
      try {
        while (true) {
          const { payload } = yield take(INPUT_STEP.SUBMIT);
          if (payload) {
            if (payload === 'fail') {
              yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
            } else {
              yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
            }
          }
        }
      } catch (e) {
        CommonLog.lError(e, {at:'input step doing'});
      }
    },
  }
}
const InputStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsInputStep extends ClsBaseStep {
  constructor(...args) {
    super(...args);
    this._statusTasks=inputStepStatusMixin(this._statusTasks);
  }
};
export default InputStepMixin;
