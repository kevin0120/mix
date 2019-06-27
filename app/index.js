import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';

import blue from '@material-ui/core/colors/blue';

import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';

import { CreateDailyLogger, Info } from './logger';

import  i18n  from './i18n'; // 初始化i18n配置

import Root from './containers/Root';
import { configureStore, history } from './store/configureStore';
import rootSaga from './modules/indexSaga';
import './app.global.css';
import { primaryColor } from './common/jss/material-react-pro';

const store = configureStore(); // 创建默认state
store.runSaga(rootSaga);

const theme = createMuiTheme({
  palette: {
    primary: blue,
    secondary: {
      main: primaryColor
    },
    status: {
      danger: 'orange'
    }
  },
  typography: {
    useNextVariants: true,
    // Use the system font instead of the default Roboto font.
    fontFamily: [
      'Noto Sans SC',
      'sans-serif'
    ].join(','),
    // fontWeightRegular: 'bold',
    button: {
      fontSize: '20px'
    }
  }
});

CreateDailyLogger(); // 创建日志对象,永远成功

Info('程序已启动...');

render(
  <AppContainer>
    <MuiThemeProvider theme={theme}>
      <Root store={store} history={history} />
    </MuiThemeProvider>
  </AppContainer>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./containers/Root', () => {
    // eslint-disable-next-line global-require
    const NextRoot = require('./containers/Root').default;
    render(
      <AppContainer>
        <MuiThemeProvider theme={theme}>
          <NextRoot store={store} history={history} />
        </MuiThemeProvider>
      </AppContainer>,
      document.getElementById('root')
    );
  });
}
