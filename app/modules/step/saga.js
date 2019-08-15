import { cancel, fork, join, put, take, select, call } from 'redux-saga/effects';
import STEP_STATUS from './model';
import stepTypes from './stepTypes';
import { orderActions, ORDER } from '../order/action';
import { orderStepUpdateApi } from '../../api/order';
import { stepStatus, workingOrder, workingStep } from '../order/selector';
import {CommonLog} from '../../common/utils';

let statusTask = null;

function invalidStepStatus(stepType, status) {
  if (!stepType) {
    throw Error(`invalid stepType ${stepType}`);
  }
  if (!status) {
    throw Error(`trying to invalid status ${status} of ${stepType}`);
  }
  throw Error(`step type ${stepType}  has empty status ${status}`);
}

function* runStep(stepType) {
  try {
    while (true) {
      const { status } = yield take(ORDER.STEP.STATUS);
      if (statusTask) {
        yield cancel(statusTask);
      }
      const id = yield select(s => workingStep(workingOrder(s.order)).id);
      yield call(orderStepUpdateApi, id, status);
      statusTask = yield fork(stepTypes?.[stepType]?.[status] ||
        (() => invalidStepStatus(stepType, status)), ORDER, orderActions);

    }
  } catch (e) {
    CommonLog.lError(e,{
      at:'runStep'
    });
  }
}

export default function* (stepType) {
  try {
    const step = yield fork(runStep, stepType);
    yield put(orderActions.stepStatus(STEP_STATUS.ENTERING));
    yield join(step);
  } catch (e) {
    CommonLog.lError(e,{
      at:'step root'
    });  }
}

