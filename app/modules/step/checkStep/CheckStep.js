// @flow
import { put, take, race } from 'redux-saga/effects';
import {STEP_STATUS} from '../constants';
import { CHECK_STEP } from './action';
import { CommonLog } from '../../../common/utils';
import type { IWorkStep } from '../interface/IWorkStep';
import {orderActions} from '../../order/action';


function checkStepStatusMixin(superTasks){
  return{
    ...superTasks,
    *[STEP_STATUS.DOING]() {
      try {
        const { submit, cancel } = yield race({
          submit: take(CHECK_STEP.SUBMIT),
          cancel: take(CHECK_STEP.CANCEL)
        });
        if (submit) {
          yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
        }
        if (cancel) {
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
        }
      } catch (e) {
        CommonLog.lError(`CheckStep DOING Error: ${e}`);
      }
    },
  }
}

const CheckStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsCheckStep extends ClsBaseStep {
  constructor(...args) {
    super(...args);
    this._statusTasks=checkStepStatusMixin(this._statusTasks);
  }
};

export default CheckStepMixin;
