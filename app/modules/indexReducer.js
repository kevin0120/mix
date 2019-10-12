// @flow

import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import setting from './setting/reducer';
import Configuration from './connections/reducer';
import users from './user/reducer';
import workMode from './workmode/reducer';
import operations from './operation/reducer';
import resultDiag from './resultDiag/reducer';
import ongoingOperation from './ongoingOperation/reducer';
import healthCheckResults from './healthzCheck/reducer';
import timeline from './timeline/reducer';
import operationViewer from './operationViewer/reducer';
import logo from './logo/reducer';
import network from './network/reducer';
import battery from './battery/reducer';
import healthz from './healthz/reducer';
import order from './order/reducer';
import dialog from './dialog/reducer';
import Notifier from './Notifier/reducer';
import loading from './loading/reducer';

export default function createRootReducer(history: {}) {
  const routerReducer = connectRouter(history)(() => {});
  return connectRouter(history)(
    combineReducers({
      router: routerReducer,
      setting,
      ...Configuration,
      users,
      workMode,
      operations,
      healthCheckResults,
      resultDiag,
      ongoingOperation,
      timeline,
      operationViewer,
      logo,
      network,
      battery,
      // tools,
      order,
      dialog,
      Notifier,
      healthz,
      loading
    })
  );
}
