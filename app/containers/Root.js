// @flow
import React, { Component } from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import type { Store } from '../reducers/types';
import Routes from '../Routes';
import { listenToNewCar } from '../actions/scannerDevice';
import { setCardAuthListener } from '../actions/cardAuth';
import { systemInit } from '../actions/sysInit';

type Props = {
  store: Store,
  history: {}
};

export default class Root extends Component<Props> {
  render() {
    const { store, history } = this.props;
    listenToNewCar(store.dispatch, store.getState);
    setCardAuthListener(store.dispatch);
    const connectionInfo = store.getState().setting.page.odooConnection;
    store.dispatch(
      systemInit(
        connectionInfo.odooUrl.value,
        connectionInfo.aiisUrl.value,
        connectionInfo.hmiSn.value,
        store.dispatch
      )
    ); // 初始化获取默认值
    return (
      <Provider store={store}>
        <ConnectedRouter history={history}>
          <Routes />
        </ConnectedRouter>
      </Provider>
    );
  }
}
