// @flow
import { put, take, race } from 'redux-saga/effects';
import {STEP_STATUS} from '../constants';
import { CHECK_STEP } from './action';
import { CommonLog } from '../../../common/utils';
import type { IWorkStep } from '../interface/IWorkStep';
import {orderActions} from '../../order/action';

const CheckStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsCheckStep extends ClsBaseStep {
  _statusTasks = {
    *[STEP_STATUS.READY](){
      try {
        yield put(orderActions.stepStatus(this, STEP_STATUS.ENTERING));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    *[STEP_STATUS.ENTERING]() {
      try {
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        CommonLog.lError(`CheckStep Entering Error: ${e}`);
      }
    },
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
    *[STEP_STATUS.FINISHED]() {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        console.error(e);
      }
    },
    *[STEP_STATUS.FAIL]() {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(`CheckStep FAIL Error: ${e}`);
      }
    }
  };
};
export default CheckStepMixin;
