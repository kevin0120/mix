import { all, call, put } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import { getDevicesByType } from '../../deviceManager/devices';
import { deviceType } from '../../deviceManager/constants';
import { orderActions } from '../../order/action';
import { CommonLog } from '../../../common/utils';
import stepActions from '../actions';

function* enteringState() {
  try {
    this._scanners = getDevicesByType(deviceType.scanner);
    const scanners = this._scanners;
    yield all(
      scanners.map(s => {
        this._listeners.push(
          s.addListener(() => true, stepActions.getValue)
        );
        return call(s.Enable);
      })
    );
    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
  } catch (e) {
    CommonLog.lError(e);
  }
}

export const scannerStepStatusMixin = (superTasks) => ({
  ...superTasks,
  [STEP_STATUS.ENTERING]: enteringState
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
