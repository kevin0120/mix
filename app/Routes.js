import React from 'react';
import { Switch, Route } from 'react-router';

// import { AnimatedSwitch } from 'react-router-transition';

import App from './containers/App';

import indexRoutes, { routeConfigs } from './routes/index';

import Layout from './components/Layout/layout';

const lodash = require('lodash');

const i = lodash.concat(routeConfigs, indexRoutes);

export default class Routes extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showLayout: false
    };
  }

  dummyFunction = (RouteConfig) => {
    const {showLayout} = this.state;
    if (RouteConfig.showLayout !== showLayout) {
      this.setState({ showLayout: RouteConfig.showLayout });
    }

    return <RouteConfig.main />;
  };

  render() {
    const {showLayout} = this.state;
    return (
      <App>
        <Switch>
          {i.map(RouteConfig => (
            <Route
              key={RouteConfig.url}
              exact
              path={RouteConfig.url}
              render={() => this.dummyFunction(RouteConfig)}
            />
          ))}
        </Switch>
        <Layout shouldRender={showLayout} />
      </App>
    );
  }
}
