// @flow

import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import setting from './setting/reducer';
import Configuration from './connections/reducer';
import users from './user/reducer';
import workMode from './workmode/reducer';
import operationViewer from './operationViewer/reducer';
import logo from './logo/reducer';
import network from './network/reducer';
import battery from './battery/reducer';
import healthz from './healthz/reducer';
import order from './order/reducer';
import dialog from './dialog/reducer';
import Notifier from './Notifier/reducer';
import loading from './loading/reducer';
import systemInfo from './systemInfo/reducer';
import workCenterMode from './workCenterMode/reducer';
import reworkPattern from './reworkPattern/reducer';
import { rootReducer } from './indexModels';

export default function createRootReducer(history: {}) {
  const routerReducer = connectRouter(history)(() => {
  });
  return connectRouter(history)(
    combineReducers({
      modules: rootReducer,
      router: routerReducer,
      setting,
      ...Configuration,
      users,
      workMode,
      workCenterMode,
      operationViewer,
      logo,
      network,
      battery,
      // tools,
      order,
      dialog,
      Notifier,
      healthz,
      loading,
      systemInfo,
      reworkPattern
    })
  );
}
