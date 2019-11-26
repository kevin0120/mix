import { IO } from './constants';

export default {
  setModule: (io) => ({
    type: IO.SET_MODULE,
    io
  }),
  set: (group, status) => ({
    type: IO.SET,
    group,
    status
  }),
  addListener: (inputType, action) => ({
    type: IO.ADD_LISTENER,
    inputType, action
  }),
  removeListener: (inputType) => ({
    type: IO.ADD_LISTENER,
    inputType
  }),
  setPortsConfig: (ioPorts) => ({
    type: IO.SET_PORT_CONFIG,
    ioPorts
  }),
  portsConfigChange: (ioPorts) => ({
    type: IO.PORT_CONFIG_CHANGE,
    ioPorts
  }),
  test: (port) => ({
    type: IO.TEST,
    port
  })
};
