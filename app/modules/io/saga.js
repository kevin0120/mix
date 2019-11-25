// @flow
import type { Saga } from 'redux-saga';
import { takeEvery, select, call, all, put } from 'redux-saga/effects';
import { IO } from './constants';
import { newDevice } from '../deviceManager/devices';
import { CommonLog } from '../../common/utils';
import { ioDirection, ioTriggerMode } from '../device/io/constants';
import { deviceType } from '../deviceManager/constants';
import ioActions from './action';

const eSetting = require('electron-settings');

let defaultIOModule = null;

const listeners = {};

export default function* root(): Saga<void> {
  try {
    defaultIOModule = newDevice(
      deviceType.io,
      'defaultIO',
      'defaultIO',
      { input_num: 8, output_num: 8 },
      {},
      []
    );
    yield put(ioActions.setModule(defaultIOModule));
    yield takeEvery(IO.SET, setIOOutput);
    yield takeEvery(IO.ADD_LISTENER, bindIOListeners);
    yield takeEvery(IO.SET_PORT_CONFIG, setPortsConfig);
    yield takeEvery(IO.TEST, testPort);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* setPortsConfig(action) {
  try {
    if (!action) {
      return;
    }
    const { ioPorts } = action;
    console.log(ioPorts);
    const {[ioDirection.input]:inputs} = ioPorts;
    const {[ioDirection.output]:outputs} = ioPorts;

    eSetting.set('devices.io.inputs', inputs);
    eSetting.set('devices.io.inputs', outputs);
    yield put(ioActions.portsConfigChange(ioPorts));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* testPort(action) {
  try {
    if (!action) {
      return;
    }
    const { port } = action;

  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindIOListeners(action) {
  try {
    if (!action || !defaultIOModule) {
      return;
    }
    const { inputType, action: act } = action;
    const { ioPorts } = yield select(s => s.io);
    const port = defaultIOModule.getPort(ioDirection.output, ioPorts[ioDirection.input][inputType]);
    if (listeners[inputType]) {
      defaultIOModule.removeListener(listeners[inputType]);
    }
    listeners[inputType] = defaultIOModule.addListener(input =>
      port === input.port &&
      ioTriggerMode.falling === input.triggerMode, act);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* setIOOutput(action) {
  try {
    if (!action || !defaultIOModule) {
      return;
    }
    const { group, status } = action;
    const { ioPorts } = yield select(s => s.io);
    console.log(group);
    const ports = group.map(o => defaultIOModule && defaultIOModule.getPort(ioDirection.output, ioPorts[ioDirection.input][o]));
    yield all(ports.map(p => call(defaultIOModule && defaultIOModule.setIO(p, status))));
  } catch (e) {
    CommonLog.lError(e);
  }
}


