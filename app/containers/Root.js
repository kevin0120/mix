// @flow
import React, { Component } from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import type { Store } from '../modules/indexReducer';
import Routes from '../Routes';
import { listenToNewCar } from '../modules/scanner/action';
import { systemInit } from '../modules/systemInit/action';

type Props = {
  store: Store,
  history: {}
};

export default class Root extends Component<Props> {
  render() {
    const { store, history } = this.props;
    listenToNewCar(store.dispatch, store.getState);
    store.dispatch(systemInit()); // 初始化
    return (
      <Provider store={store}>
        <ConnectedRouter history={history}>
          <Routes />
        </ConnectedRouter>
      </Provider>
    );
  }
}
