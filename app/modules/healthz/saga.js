import { takeEvery } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { HEALTHZ } from './action';

import { updateDeviceStatus } from '../external/device';

export default function* healthz() {
  try {
    yield takeEvery(HEALTHZ.UPDATE, updateDeviceStatus);
  } catch (e) {
    CommonLog.lError(e);
  }
}