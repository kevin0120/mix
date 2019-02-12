import React from 'react';
import { Switch, Route } from 'react-router';

// import { AnimatedSwitch } from 'react-router-transition';

import App from './containers/App';

import indexRoutes, { routeConfigs } from './routes/index';

import Layout from './components/Layout/layout';

const lodash = require('lodash');

const i = lodash.concat(routeConfigs, indexRoutes);

export default class Routes extends React.Component {

  renderRoute = RouteConfig => {
    return (
      <React.Fragment>
        <Layout shouldRender={RouteConfig.showLayout}>
          <RouteConfig.main/>
        </Layout>
      </React.Fragment>
    );
  };

  render() {
    return (
      <App>
        <Switch>
          {i.map(RouteConfig => (
            <Route
              key={RouteConfig.url}
              exact
              path={RouteConfig.url}
              render={() => {
                return this.renderRoute(RouteConfig);
              }}
            />
          ))}
        </Switch>
      </App>
    );
  }
}
