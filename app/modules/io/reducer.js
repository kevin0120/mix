import { IO, ioOutputGroups } from './constants';
import { genReducers } from '../util';
import { ioDirection } from '../device/io/constants';

const eSetting = require('electron-settings');

const ioConfig = eSetting.get('devices.io');

const reducers = {
  // [IO.SET_PORT_CONFIG]: (state, { output, portIdx }) => ({
  //   ...state,
  //   [output]: portIdx
  // }),
  [IO.SET_MODULE]: (state, { io }) => ({
    ...state,
    ioModule: io
  }),
  [IO.PORT_CONFIG_CHANGE]: (state, { ioPorts }) => ({
    ...state,
    ioPorts
  }),
};

const initState = {
  enabled: ioConfig.enable,
  ioPorts: {
    [ioDirection.output]: ioConfig.outputs || {},
    [ioDirection.input]: ioConfig.inputs || {}
  },
  ioOutputGroups,
  ioModule: null,
  testStatus: []
};

export default genReducers(reducers, initState);
