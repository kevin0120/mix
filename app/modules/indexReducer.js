// @flow

import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';
import setting from './setting/reducer';
import notify from './notification/reducer';
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
import tools from './tools/reducer';
import order from './order/reducer';
import dialog from './dialog/reducer';

export type StateType = {
  +notify: Object,
  +setting: Object
};

export type Action = {
  +type: string
};

export type GetState = () => StateType;

export type Dispatch = ReduxDispatch<Action>;

export type Store = ReduxStore<GetState, Action>;


export default function createRootReducer(history: {}) {
  const routerReducer = connectRouter(history)(() => {});
  return connectRouter(history)(
    combineReducers({
      router: routerReducer,
      setting,
      notify,
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
      tools,
      order,
      dialog
    })
  );
}
