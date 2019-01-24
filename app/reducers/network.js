import { NETWORK } from '../actions/actionTypes';
import configs from '../shared/config';
import { genReducers } from './utils';

const initNetwork = {
  config: configs.page.network,
  ssidList: [],
  connecting: false
};

const reducers = {
  [NETWORK.SET_OK]: networkSetOk,
  [NETWORK.SET]: networkSet,
  [NETWORK.SET_FAIL]: networkSetFail,
  [NETWORK.SCAN_OK]: networkScanOk,
  [NETWORK.CHECK_OK]: networkCheckOk,
  [NETWORK.SIGNAL_OK]: networkSignalOk,
};

function networkSetOk(state, action) {
  return {
    ...state,
    config: action.data,
    connecting: false
  };
}

function networkSet(state, action) {
  return {
    ...state,
    connecting: true
  };
}

function networkSetFail(state, action) {
  return {
    ...state,
    connecting: false
  };
}

function networkScanOk(state, action) {
  return {
    ...state,
    ssidList: action.data
  };
}

function networkCheckOk(state, action) {
  return {
    ...state,
    config: action.data
  };
}

function networkSignalOk(state, action) {
  return {
    ...state,
    ssid: action.ssid,
    signal: action.signal
  };
}

export default genReducers(reducers, initNetwork);

// export default function network(
//   state: object = initNetwork,
//   action
// ) {
//   switch (action.type) {
//     case NETWORK.SET_OK:
//       return {
//         ...state,
//         config: action.data,
//         connecting: false
//       };
//     case NETWORK.SET:
//       return {
//         ...state,
//         connecting: true
//       };
//     case NETWORK.SET_FAIL:
//       return {
//         ...state,
//         connecting: false
//       };
//     case NETWORK.SCAN_OK:
//       return {
//         ...state,
//         ssidList: action.data
//       };
//     case NETWORK.CHECK_OK:
//       return {
//         ...state,
//         config: action.data
//       };
//     case NETWORK.SIGNAL_OK:
//       return {
//         ...state,
//         ssid: action.ssid,
//         signal: action.signal
//       };
//     default:
//       return state;
//   }
// }
