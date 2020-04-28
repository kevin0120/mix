import { model, watch, write } from 'saga-modux';
import { all, call, delay, put, select } from 'redux-saga/effects';
import { ioDirection, ioTriggerMode } from '../device/io/constants';
import { ioOutputGroups } from './constants';
import { newDevice } from '../deviceManager/devices';
import { deviceType } from '../deviceManager/constants';
import { CommonLog } from '../../common/utils';

const eSetting = require('electron-settings');

const ioConfig = eSetting.get('devices.io');

let defaultIOModule = null;

const listeners = {};

const ioFunctions = {};

export default model('io', {
  * root() {
    try {
      yield call(this._initIO);
      yield call(watch([
        this.bindIOListeners,
        this.removeIOListener,
        this.setPortsConfig,
        this.testPort,
        this.setIOOutput
      ]));
    } catch (e) {
      CommonLog.lError(e, { at: 'io model' });
    }
  },

  * _initIO() {
    try {
      defaultIOModule = newDevice(
        deviceType.io,
        'defaultIO',
        eSetting.get('devices.io.sn'),
        eSetting.get('devices.io.config'),
        {},
        []
      );
      // TODO set init io status
      yield call(this.setIOOutput, { group: ioOutputGroups.ready, status: true });
      yield write(this, (s) => ({
        ...s,
        ioModule: defaultIOModule
      }));
    } catch (e) {
      console.error(e);
    }
  },


  * setPortsConfig({ ioPorts }) {
    try {
      if (!ioPorts) {
        return;
      }
      const { [ioDirection.input]: inputs } = ioPorts;
      const { [ioDirection.output]: outputs } = ioPorts;

      const { ioPorts: prevIoPorts } = yield select(this.select);
      eSetting.set('devices.io.inputs', inputs, {
        prettify: true
      });
      eSetting.set('devices.io.inputs', outputs, {
        prettify: true
      });
      yield write(this, s => ({
        ...s,
        ioPorts
      }));

      const {
        [ioDirection.input]: prevInputs,
        [ioDirection.output]: prevOutputs
      } = prevIoPorts || {};
      const effects = [];
      if (prevInputs) {
        Object.keys(prevInputs).forEach((k) => {
          if (prevInputs[k] !== inputs[k]) {
            effects.push(put(this.removeIOListener({ inputType: k })));
          }
        });
        Object.keys(inputs).forEach((k) => {
          if (prevInputs[k] !== inputs[k]) {
            effects.push(put(this.addListener({
              inputType: k,
              action: ioFunctions[k]
            })));
          }
        });
      }

      const { ioOutStatus } = yield select(this.select);
      if (prevOutputs) {

        Object.keys(prevOutputs).forEach((k) => {
          if (prevOutputs[k] !== outputs[k]) {
            effects.push(put(this.setIOOutput({ group: [k], status: ioOutStatus[prevOutputs[k]] })));
          }
        });
        Object.keys(outputs).forEach((k) => {
          if (prevOutputs[k] !== outputs[k]) {
            effects.push(put(this.setIOOutput({ group: [k], status: ioOutStatus[outputs[k]] })));
          }
        });
      }

      yield all(effects);
    } catch (e) {
      CommonLog.lError(e);
    }
  },

  * testPort({ port } = {}) {
    try {
      if (!port || !defaultIOModule) {
        return;
      }
      const { direction, idx } = port;
      const ioPort = defaultIOModule.getPort(direction, idx);
      yield call(defaultIOModule.openIO, ioPort);
      yield delay(3000);
      yield call(defaultIOModule.closeIO, ioPort);
    } catch (e) {
      CommonLog.lError(e);
    }
  },

  removeIOListener({ inputType } = {}) {
    try {
      if (!inputType || !defaultIOModule) {
        return;
      }
      if (listeners[inputType]) {
        defaultIOModule.removeListener(listeners[inputType]);
        delete listeners[inputType];
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  },

  * bindIOListeners({ inputType, action } = {}) {
    try {
      if (!action || !inputType || !defaultIOModule) {
        return;
      }
      const { ioPorts } = yield select(this.select);
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
        action
      );
    } catch (e) {
      CommonLog.lError(e);
    }
  },

  * setIOOutput({ group, status } = {}) {
    try {
      if (!group || !defaultIOModule) {
        return;
      }
      const { on, off } = group;
      const { ioPorts } = yield select(this.select);
      console.log(on, status, ioPorts);

      const onPorts = (on || []).map(
        o =>
          defaultIOModule &&
          defaultIOModule.getPort(
            ioDirection.output,
            ioPorts[ioDirection.output][o]
          )
      );
      const offPorts = (off || []).map(
        o =>
          defaultIOModule &&
          defaultIOModule.getPort(
            ioDirection.output,
            ioPorts[ioDirection.output][o]
          )
      );

      if (status) {
        yield all(offPorts.map(p => call(
          (defaultIOModule && defaultIOModule.setIO) || (() => {
          }), p, false)
        ));
      }
      yield all(onPorts.map(p => call(
        (defaultIOModule && defaultIOModule.setIO) || (() => {
        }), p, status)
      ));

    } catch (e) {
      CommonLog.lError(e);
    }
  }

}, {
  enabled: ioConfig.enable,
  ioPorts: {
    [ioDirection.output]: ioConfig.outputs || {},
    [ioDirection.input]: ioConfig.inputs || {}
  },
  ioOutputGroups,
  ioModule: null,
  testStatus: [],
  ioOutStatus: [...new Array(ioConfig.config.output_num)].map(() => 0)
});
