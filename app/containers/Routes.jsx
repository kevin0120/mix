import React from 'react';
import { Route } from 'react-router';
import filterRoutesByConfig from './pages';
import Layout from '../components/Layout/layout';
import Pages from './layouts/Pages';

class Routes extends React.Component {

  render() {
    const { pagesConfig } = this.props;
    const pages = filterRoutesByConfig(pagesConfig);
    console.log(pages);
    return (
      <React.Fragment>
        {pages.map(page => (
            <Route
              key={page.url}
              exact
              path={page.url}
              component={page.main}
            />
        ))}
        <Route key="/app" path="/app" render={() => <Layout pages={pages}/>}/>
      </React.Fragment>
    );
  }
}

export default Routes;
