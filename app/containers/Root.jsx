// @flow

import React, { Component } from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import blue from '@material-ui/core/colors/blue';
import { MuiThemeProvider } from '@material-ui/core/styles';
import type { Store } from '../modules/indexReducer';
import Routes from './Routes';
import theme from '../common/theme';

type Props = {
  store: Store,
  history: {}
};

class Root extends Component<Props> {
  render() {
    const { store, history } = this.props;
    const { pages } = store.getState().setting;
    // listenToNewCar(store.dispatch, store.getState);
    // store.dispatch(systemInit()); // 初始化
    return (
      <MuiThemeProvider theme={theme}>
        <Provider store={store}>
          <ConnectedRouter history={history}>
            <div style={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Routes pagesConfig={pages}/>
            </div>
          </ConnectedRouter>
        </Provider>
      </MuiThemeProvider>
    );
  }
}

export default Root;
