// @flow
import React from 'react';
import { Route } from 'react-router';
import pages from './pages';

type Props = {
  pagesConfig: {}
};

function renderRoute(R, subRouteList) {
  return R ?
    <Route
      key={R.url}
      exact={R.exact || false}
      path={R.url}
      render={() => R.component ?
        <R.component self={R} childRoutes={subRouteList[0]}>
          {subRouteList[1]}
        </R.component> : null}
    /> : null;
}

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
    const childRoutes=parseRouteTree(route, route.url, filter[route.name]);
    return renderRoute(route, childRoutes);
  });
  return [routeList, renderedRoute];
}

export default function Routes(props: Props) {
  const { pagesConfig } = props;
  const RoutesParsed = parseRouteTree(pages, '', pagesConfig);
  return (
    <React.Fragment>
      {RoutesParsed[1]}
    </React.Fragment>
  );

}
