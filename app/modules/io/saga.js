// @flow
import type { Saga } from 'redux-saga';
import { takeEvery, select, call, all, put, delay } from 'redux-saga/effects';
import { IO, ioInputs, ioOutputGroups } from './constants';
import { newDevice } from '../deviceManager/devices';
import { CommonLog } from '../../common/utils';
import { ioDirection, ioTriggerMode } from '../device/io/constants';
import { deviceType } from '../deviceManager/constants';
import ioActions from './action';

const eSetting = require('electron-settings');

let defaultIOModule = null;

const listeners = {};

const ioFunctions = {
  // [ioInputs.resetKey]:
};

export default function* root(): Saga<void> {
  try {
    yield call(initIO);
    yield put(ioActions.setModule(defaultIOModule));
    yield takeEvery(IO.SET, setIOOutput);
    yield takeEvery(IO.ADD_LISTENER, bindIOListeners);
    yield takeEvery(IO.REMOVE_LISTENER, removeIOListener);
    yield takeEvery(IO.SET_PORT_CONFIG, setPortsConfig);
    yield takeEvery(IO.TEST, testPort);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* initIO() {
  try {
    defaultIOModule = newDevice(
      deviceType.io,
      'defaultIO',
      eSetting.get('devices.io.sn'),
      { input_num: 8, output_num: 8 },
      {},
      []
    );
    // TODO set init io status

    yield put(ioActions.set(ioOutputGroups.ready, true));
  } catch (e) {
    console.error(e);
  }
}

function* setPortsConfig(action) {
  try {
    if (!action) {
      return;
    }
    const { ioPorts } = action;
    const { [ioDirection.input]: inputs } = ioPorts;
    const { [ioDirection.output]: outputs } = ioPorts;

    const { ioPorts: prevIoPorts } = yield select(s => s.io);
    eSetting.set('devices.io.inputs', inputs);
    eSetting.set('devices.io.inputs', outputs);
    yield put(ioActions.portsConfigChange(ioPorts));

    const {
      [ioDirection.input]: prevInputs,
      [ioDirection.output]: prevOutputs
    } = prevIoPorts || {};
    const effects = [];
    console.log(prevInputs, prevOutputs, inputs, outputs);
    if (prevInputs) {
      Object.keys(prevInputs).forEach((k) => {
        if (prevInputs[k] !== inputs[k]) {
          effects.push(put(ioActions.removeListener(k)));
        }
      });
      Object.keys(inputs).forEach((k) => {
        if (prevInputs[k] !== inputs[k]) {
          effects.push(put(ioActions.addListener(k, ioFunctions[k])));
        }
      });
    }

    const { ioOutStatus } = yield select(s => s.io);
    if (prevOutputs) {

      Object.keys(prevOutputs).forEach((k) => {
        if (prevOutputs[k] !== outputs[k]) {
          effects.push(put(ioActions.set([k], ioOutStatus[prevOutputs[k]])));
        }
      });
      Object.keys(outputs).forEach((k) => {
        if (prevOutputs[k] !== outputs[k]) {
          effects.push(put(ioActions.set([k], ioOutStatus[outputs[k]])));
        }
      });
    }

    yield all(effects);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* testPort(action) {
  try {
    if (!action || !defaultIOModule) {
      return;
    }
    const {
      port: { direction, idx }
    } = action;
    const port = defaultIOModule.getPort(direction, idx);
    yield call(defaultIOModule.openIO, port);
    yield delay(3000);
    yield call(defaultIOModule.closeIO, port);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function removeIOListener(action) {
  try {
    if (!action || !defaultIOModule) {
      return;
    }
    const { inputType } = action;
    if (listeners[inputType]) {
      defaultIOModule.removeListener(listeners[inputType]);
      delete listeners[inputType];
    }
    console.log(listeners);
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
    const port = defaultIOModule.getPort(
      ioDirection.output,
      ioPorts[ioDirection.input][inputType]
    );
    if (listeners[inputType]) {
      defaultIOModule.removeListener(listeners[inputType]);
      delete listeners[inputType];
    }
    listeners[inputType] = defaultIOModule.addListener(
      input =>
        port === input.port && ioTriggerMode.falling === input.triggerMode,
      act
    );
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
    console.log(group, status);


    const ports = group.map(
      o =>
        defaultIOModule &&
        defaultIOModule.getPort(
          ioDirection.output,
          ioPorts[ioDirection.output][o]
        )
    );

    console.log(ports);
    if (status) {
      const otherPorts = defaultIOModule.ports
        .filter(p => p.direction === ioDirection.output)
        .filter(p => (ports || []).every(sp => sp.idx !== p.idx));
      yield all(otherPorts.map(p => call(
        (defaultIOModule && defaultIOModule.setIO) || (() => {
        }), p, false)
      ));
    }
    yield all(ports.map(p => call(
      (defaultIOModule && defaultIOModule.setIO) || (() => {
      }), p, status)
    ));

  } catch (e) {
    CommonLog.lError(e);
  }
}

