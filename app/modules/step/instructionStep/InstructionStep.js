// @flow
import { put, take } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import { INSTRUCTION_STEP } from './action';
import type { IWorkStep } from '../interface/IWorkStep';
import { CommonLog } from '../../../common/utils';
import {orderActions} from '../../order/action';

const instructionStepStatusMixin=(superTasks)=>({
  ...superTasks,
  * [STEP_STATUS.DOING]() {
    try {
      yield take(INSTRUCTION_STEP.SUBMIT);
      yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
    } catch (e) {
      console.error(e);
    }
  },
});

const InstructionStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsInstructionStep extends ClsBaseStep {
  constructor(...args) {
    super(...args);
    this._statusTasks=instructionStepStatusMixin(this._statusTasks);
  }
};
export default InstructionStepMixin;
