import React from 'react';
import { Switch, Route } from 'react-router';
import App from './containers/App';

import indexRoutes, { routeConfigs } from './routes/index';

const lodash = require('lodash');

const i = lodash.concat(routeConfigs, indexRoutes);

export default () => (
  <App>
    <Switch>
      {i.map(route => (
        <Route key={route.url} exact path={route.url} component={route.main} />
      ))}
    </Switch>
  </App>
);
