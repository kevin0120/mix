import { actionChannel, call, put, select, take } from 'redux-saga/effects';
import { orderActions } from '../order/action';
import { STEP_ACTIONS, STEP_STATUS } from './constants';
import { CommonLog } from '../../common/utils';
import dialogActions from '../dialog/action';
import actions from './actions';
import { workModes } from '../workCenterMode/constants';

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
    const stepChannel = yield actionChannel([
      STEP_ACTIONS.INPUT,
      STEP_ACTIONS.SUBMIT
    ]);
    while (true) {
      const action = yield take(stepChannel);
      switch (action.type) {
        case STEP_ACTIONS.INPUT:
          yield call(this.updateData, d => ({
            ...(d || {}),
            result: action?.input?.data,
            timeLine: [
              {
                title: action?.input?.source,
                color: 'info',
                footerTitle:
                  action &&
                  action.input &&
                  action.input.time.toLocaleString(),
                body: String(action?.input?.data)
              },
              ...(d?.timeLine || [])
            ]
          }));
          break;
        case STEP_ACTIONS.SUBMIT: {
          const { result } = action;
          yield call(this.updateData, d => ({
            ...(d || {}),
            result
          }));
          yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
          break;
        }
        default:
          break;
      }
    }
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
    const { error } = config;
    const { workCenterMode } = yield select();
    const isNormal = workCenterMode === workModes.normWorkCenterMode;
    let buttons = [
      {
        label: 'Common.Close',
        color: 'danger',
        action: actions.confirmFail()
      }
    ];
    if (isNormal) {
      buttons = buttons.concat([
        {
          label: 'Order.Next',
          color: 'warning',
          action: actions.confirmFail()
        }
      ]);
    }
    yield put(
      dialogActions.dialogShow({
        buttons,
        title: `工步失败：${this._code}`,
        content: `${error || this.failureMsg}`
      })
    );
    yield take(STEP_ACTIONS.CONFIRM_FAIL);
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
