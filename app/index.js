import { hot } from 'react-hot-loader/root';

import React from 'react';
import { render } from 'react-dom';
import { CreateDailyLogger, Info } from './logger';
import Root from './containers/Root';
import { configureStore, history } from './store/configureStore';
import rootSaga from './modules/indexSaga';
import './app.global.css';
import { loadLocales } from './locales';

const store = configureStore(); // 创建默认state
store.runSaga(rootSaga);


CreateDailyLogger(); // 创建日志对象,永远成功

Info('程序已启动...');

loadLocales();

const HotRoot = hot(Root);

render(
  <HotRoot store={store} history={history}/>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./containers/Root', () => {
    // eslint-disable-next-line global-require
    const NextRoot = require('./containers/Root').default;
    const HotNextRoot = hot(NextRoot);
    render(
      <HotNextRoot store={store} history={history}/>,
      document.getElementById('root')
    );
  });
}
