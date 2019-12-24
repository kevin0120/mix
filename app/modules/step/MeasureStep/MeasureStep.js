// @flow
import { call, put, race, take } from 'redux-saga/effects';
import { STEP_STATUS,STEP_ACTIONS } from '../constants';
import { CommonLog } from '../../../common/utils';
import type { IWorkStep } from '../interface/IWorkStep';
import { orderActions } from '../../order/action';


function checkStepStatusMixin(superTasks) {
  return {
    ...superTasks,
    * [STEP_STATUS.DOING]() {
      try {
        const { submit } = yield race({
          submit: take(STEP_ACTIONS.SUBMIT)
        });
        if (submit) {
          yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
          const { result } = submit;
          yield call(
            this.updateData, (data) => ({
              ...data,
              result
            })
          );
        }
      } catch (e) {
        CommonLog.lError(`CheckStep DOING Error: ${e}`);
      }
    }
  };
}

const CheckStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsCheckStep extends ClsBaseStep {
  constructor(...args) {
    super(...args);
    this._statusTasks = checkStepStatusMixin(this._statusTasks);
  }
};

export default CheckStepMixin;
