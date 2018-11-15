// @flow
import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import setting from './setting';
import notify from './notify';
import Configuration from './configuration'


export default function createRootReducer(history: {}) {
  const routerReducer = connectRouter(history)(() => {});

  return connectRouter(history)(
    combineReducers({ router: routerReducer, setting, notify, ...Configuration })
  );
}
