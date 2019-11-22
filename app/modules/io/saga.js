// @flow
import { takeEvery, select, call, all, put } from 'redux-saga/effects';
import { IO } from './constants';
import { newDevice } from '../deviceManager/devices';
import { CommonLog } from '../../common/utils';
import { ioDirection, ioTriggerMode } from '../device/io/constants';
import { deviceType } from '../deviceManager/constants';
import ioActions from './action';

let defaultIOModule = null;

const listeners = {};

export default function* root() {
  try {
    defaultIOModule = newDevice(
      deviceType.io,
      'defaultIO',
      'defaultIO',
      { input_num: 8, output_num: 8 }
    );
    yield put(ioActions.setModule(defaultIOModule));
    yield takeEvery(IO.SET, setIOOutput);
    yield takeEvery(IO.ADD_LISTENER, bindIOListeners);
    yield takeEvery(IO.SET_PORT, setPort);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* setPort(action) {
  try {
    if (!action) {
      return;
    }
    const { output, portIdx } = action;

  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindIOListeners(action) {
  try {
    if (!action) {
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
    if (!action) {
      return;
    }
    const { group, status } = action;
    const { ioPorts } = yield select(s => s.io);
    const ports = group.map(o => defaultIOModule.getPort(ioDirection.output, ioPorts[ioDirection.input][o]));
    yield all(ports.map(p => call(defaultIOModule.setIO(p, status))));
  } catch (e) {
    CommonLog.lError(e);
  }
}


