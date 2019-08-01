// @flow

import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { MuiThemeProvider } from '@material-ui/core/styles';
import SnackbarProvider from '../components/Notifier/SnackbarProvider';
import type { Store } from '../modules/indexReducer';
import Routes from './Routes';
import theme from '../common/theme';

type Props = {
  store: Store,
  history: {}
};


const Root = (props: Props) => {
  const { store, history } = props;
  const { pages } = store.getState().setting;
  // listenToNewCar(store.dispatch, store.getState);
  // store.dispatch(systemInit()); // 初始化
  return (
    <MuiThemeProvider theme={theme}>
      <Provider store={store}>
        <SnackbarProvider maxSnack={5}>
          <ConnectedRouter history={history}>
            <div style={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Routes pagesConfig={pages}/>
            </div>
          </ConnectedRouter>
        </SnackbarProvider>
      </Provider>
    </MuiThemeProvider>
  );
};

export default Root;
