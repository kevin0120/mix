import React from 'react';
import { Switch, Route } from 'react-router';
import App from './containers/App';

import indexRoutes, { routeConfigs } from './routes';

const lodash = require('lodash');

const i = lodash.concat(routeConfigs, indexRoutes);

export default () => (
  <App>
    <Switch>
      {i.map(route => (
        <Route
          key={route.url || route.path}
          exact
          path={route.url || route.path}
          component={route.component || route.main}
        />
      ))}
    </Switch>
  </App>
);
