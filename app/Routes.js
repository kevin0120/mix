import React from 'react';
import { Switch, Route } from 'react-router';
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

  render() {
    // console.log('rerender');
    return (
      <App>
        <Switch>
          {i.map(RouteConfig => (
            <Route
              key={RouteConfig.url}
              exact
              path={RouteConfig.url}
              render={() => {
                if (RouteConfig.showLayout !== this.state.showLayout) {
                  this.setState({ showLayout: RouteConfig.showLayout });
                }
                return <RouteConfig.main />;
              }}
            />
          ))}
        </Switch>
        <Layout shouldRender={this.state.showLayout} />
      </App>
    );
  }
}
