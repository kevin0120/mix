import { put, take } from 'redux-saga/effects';
import { orderActions } from '../order/action';
import { STEP_ACTIONS, STEP_STATUS } from './constants';
import { CommonLog } from '../../common/utils';

function* readyState(config) {
  try {
    yield put(orderActions.stepStatus(this, STEP_STATUS.ENTERING, config));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* enteringState(config) {
  try {
    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING, config));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* doingState(config) {
  try {
    yield take(STEP_ACTIONS.SUBMIT);
    yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED, config));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* finishState(config) {
  try {
    yield put(orderActions.finishStep(this));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* failState(config) {
  try {
    yield put(orderActions.finishStep(this));
  } catch (e) {
    console.error(e);
  }
}

export const stepStatusTasks = {
  [STEP_STATUS.READY]: readyState,
  [STEP_STATUS.ENTERING]: enteringState,
  [STEP_STATUS.DOING]: doingState,
  [STEP_STATUS.FINISHED]: finishState,
  [STEP_STATUS.FAIL]: failState
};
