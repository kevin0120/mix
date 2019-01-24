// @flow
import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import setting from './setting';
import notify from './notify';
import Configuration from './configuration';
import users from './user';
import workMode from './workmode';
import operations from './operations';
import shutDownDiag from './shutDownDiag';
import resultDiag from './reusltDiag';
import ongoingOperation from './ongoingOpeartion';
import healthCheckResults from './healthCheck';
import timeline from './timeline';
import operationViewer from './operationViewer';
import logo from './logo';
import network from './network';
import battery from './battery';

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
      shutDownDiag,
      resultDiag,
      ongoingOperation,
      timeline,
      operationViewer,
      logo,
      network,
      battery
    })
  );
}
