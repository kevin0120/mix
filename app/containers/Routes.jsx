import React from 'react';
import { Route } from 'react-router';
import filterRoutesByConfig from './pages';
import { pages } from './pages';

function renderRoute(R, subRouteList) {
  if (!R) {
    return null;
  }
  return <Route key={R.url}
                exact={R.exact || false}
                path={R.url}
                render={() => R.component ?
                  <R.component {...R.props} childRoutes={subRouteList}>
                    {subRouteList&&subRouteList.map(subRoute => renderRoute(subRoute)) || null}
                  </R.component> : null}
  />;
};

function parseRouteTree(routesObj, parentUrl, filter) {
  const routeUrls = Object.keys(routesObj).filter((key) => /^\//.test(key));
  const routeList = routeUrls.map(u => {
    const name = u.slice(1);
    return filter[name] ? {
      ...routesObj[u],
      url: (parentUrl || '') + u,
      name,
      role: filter[name] instanceof Array ? filter[name] : []
    } : null;
  });
  const renderedRoute = routeList.map((route) => {
    if (!route) {
      return null;
    }
    return renderRoute(route, parseRouteTree(route, route.url, filter[route.name])[0]);
  });
  return [routeList, renderedRoute];
}

class Routes extends React.Component {

  render() {
    const { pagesConfig } = this.props;
    console.log(pagesConfig);
    return (
      <React.Fragment>
        {parseRouteTree(pages, '', pagesConfig)[1]}
      </React.Fragment>
    );
  }
}

export default Routes;
