import { cancel, fork, join, put, take, select, call } from 'redux-saga/effects';
import STEP_STATUS from './model';
import stepTypes from './stepTypes';
import { orderActions, ORDER } from '../order/action';
import { orderStepUpdateApi } from '../../api/order';
import { stepStatus, workingOrder, workingStep } from '../order/selector';

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


export default function* (stepType) {
  try {
    const step = yield fork(
      function* runStep() {
        try {
          while (true) {
            const { status } = yield take(ORDER.STEP.STATUS);
            if (statusTask) {
              yield cancel(statusTask);
            }
            statusTask = yield fork(stepTypes?.[stepType]?.[status] || (() => invalidStepStatus(stepType, status)), ORDER, orderActions);
            const id = yield select(s => workingStep(workingOrder(s.order)).id);
            yield call(orderStepUpdateApi, id, status);
          }
        } catch (e) {
          console.error(e);
        }
      }
    );
    yield put(orderActions.stepStatus(STEP_STATUS.ENTERING));
    yield join(step);
  } catch (e) {
    console.error(e);
  }
}

