// @flow
import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import { deviceStatusApi } from '../../api/device';
import { CommonLog } from '../../common/utils';

// eslint-disable-next-line import/prefer-default-export
export function* updateDeviceStatus(): Saga<void> {
  try {
    yield call(deviceStatusApi);
  } catch (e) {
    CommonLog.lError(e);
  }
}
