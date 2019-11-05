import { takeEvery } from 'redux-saga/effects';
import { SYSTEM_INFO } from './action';
import { CommonLog } from '../../common/utils';

const eSetting = require('electron-settings');

export default function* system() {
  try {
    yield takeEvery(SYSTEM_INFO.SET_WORKCENTER, setWorkCenter);
  } catch (e) {
    CommonLog.lError(e, { at: 'systemInfo' });
  }
}


function setWorkCenter({ workcenter }) {
  eSetting.set('system.workcenter.code', workcenter);
}