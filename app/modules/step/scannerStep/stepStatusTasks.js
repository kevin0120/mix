import { all, call, put, select, take } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import { getDevicesByType } from '../../deviceManager/devices';
import { deviceType } from '../../deviceManager/constants';
import { SCANNER_STEP, scannerStepAction } from './action';
import { orderActions } from '../../order/action';
import { CommonLog } from '../../../common/utils';
import { stepData, stepPayload, workingOrder, workingStep } from '../../order/selector';

function* enteringState() {
  try {
    this._scanners = getDevicesByType(deviceType.scanner);
    const scanners = this._scanners;
    yield all(
      scanners.map(s => {
        // eslint-disable-next-line no-param-reassign
        this._listeners.push(
          s.addListener(() => true, scannerStepAction.getValue)
        );
        return call(s.Enable);
      })
    );
    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* doingState() {
  try {
    while (true) {
      const action = yield take([
        SCANNER_STEP.GET_VALUE,
        SCANNER_STEP.SUBMIT
      ]);
      const result = yield select(
        s => stepData(workingStep(workingOrder(s.order)))?.result
      );
      // const label = yield select(
      //   s => stepPayload(workingStep(workingOrder(s.order)))?.label
      // );
      switch (action.type) {
        case SCANNER_STEP.GET_VALUE:
          yield call(this.updateData, d => ({
            ...(d || {}),
            result: action?.input?.data,
            timeLine: [
              {
                title: action?.input?.name,
                color: 'info',
                footerTitle:
                  action &&
                  action.input &&
                  action.input.time.toLocaleString(),
                body: action?.input?.data
              },
              ...(d?.timeLine || [])
            ]
          }));
          break;
        case SCANNER_STEP.SUBMIT:
          // if(this._steps.length>0){
          //   yield call(this.runSubStep,this._steps[0])
          // }
          yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
          break;
        default:
          break;
      }
    }
  } catch (e) {
    CommonLog.lError(e);
  } finally {
    yield all(
      this._scanners.map(s => {
        // eslint-disable-next-line no-param-reassign
        s.dispatcher = null;
        return call(s.Disable);
      })
    );
  }
}


export const scannerStepStatusMixin = (superTasks) => ({
  ...superTasks,
  [STEP_STATUS.ENTERING]: enteringState,
  [STEP_STATUS.DOING]: doingState
});

export function onLeave() {
  (this._scanners || []).forEach(s => {
    this._listeners.forEach(l => {
      s.removeListener(l);
    });
  });
  this._scanners = [];
  this._listeners = [];
}
