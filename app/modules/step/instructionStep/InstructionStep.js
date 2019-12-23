// @flow
import { put, take } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import { INSTRUCTION_STEP } from './action';
import type { IWorkStep } from '../interface/IWorkStep';
import { CommonLog } from '../../../common/utils';
import {orderActions} from '../../order/action';

const InstructionStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsInstructionStep extends ClsBaseStep {
  _statusTasks = {
    *[STEP_STATUS.READY](){
      try {
        yield put(orderActions.stepStatus(this, STEP_STATUS.ENTERING));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    * [STEP_STATUS.ENTERING]() {
      try {
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.DOING]() {
      try {
        yield take(INSTRUCTION_STEP.SUBMIT);
        yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
      } catch (e) {
        console.error(e);
      }
    },
    * [STEP_STATUS.FINISHED]() {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        console.error(e);
      }
    }
  };
};
export default InstructionStepMixin;
