import { cancel, fork, join, put, take } from 'redux-saga/effects';
import { STEP_STATUS } from './model';
import stepTypes from './stepTypes';
import { orderActions, ORDER } from '../order/action';

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

