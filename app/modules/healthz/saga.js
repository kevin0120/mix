import { takeEvery,fork } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import healthzActions, { HEALTHZ } from './action';
import { updateDeviceStatus } from '../external/device';
import { bindOnConnectAction } from '../rush/rushHealthz';

export default function* healthz() {
  try {
    yield takeEvery(HEALTHZ.UPDATE, updateDeviceStatus);
    yield fork(bindOnConnectAction,healthzActions.update);
  } catch (e) {
    CommonLog.lError(e);
  }
}