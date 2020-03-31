import { actionChannel, call, put, take } from 'redux-saga/effects';
import { STEP_ACTIONS, STEP_STATUS } from '../constants';
import { CommonLog } from '../../../common/utils';
import { orderActions } from '../../order/action';

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
          const color = action?.input?.data ? 'info' : 'danger';
          yield call(this.updateData, d => ({
            ...(d || {}),
            result: action?.input?.data,
            timeLine: [
              {
                title: action?.input?.source,
                color,
                footerTitle:
                  action &&
                  action.input &&
                  action.input.time.toLocaleString(),
                body: action?.input?.data ? 'Pass' : 'Fail'
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

export const passfailStepStatusMixin = (superTasks) => ({
  ...superTasks,
  [STEP_STATUS.DOING]: doingState
});
