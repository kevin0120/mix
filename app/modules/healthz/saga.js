// @flow
import { takeEvery, fork } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import healthzActions, { HEALTHZ } from './action';
import { updateDeviceStatus } from '../external/device';
import { bindRushAction } from '../rush/rushHealthz';

export default function* healthz(): Saga<void> {
  try {
    yield takeEvery(HEALTHZ.UPDATE, updateDeviceStatus);
    yield fork(bindRushAction.onConnect, healthzActions.update);
  } catch (e) {
    CommonLog.lError(e);
  }
}
