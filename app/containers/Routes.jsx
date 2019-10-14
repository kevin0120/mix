// @flow

import React from 'react';
import { Route } from 'react-router';
import type { Node } from 'react';
import pages, { getContentByUrl } from './pages';
import type { tRouteObj, tRouteComponent } from './typeDef';

type Props = {
  pagesConfig: { [key: string]: Object | Array<string> }
};

type tRouteList = [
  ?Array<?tRouteObj>,
  ?Array<?tRouteComponent>
];

function renderRoute(R: ?tRouteObj, subRouteList: tRouteList): ?Node {
  return R ?
    <Route
      key={R.url}
      exact={R.exact || false}
      path={R.url}
      render={() => R.component ?
        <R.component self={R} childRoutes={subRouteList[0]} getContentByUrl={getContentByUrl}>
          {subRouteList[1]}
        </R.component> : null}
    /> : null;
}

function parseRouteTree(routesObj: tRouteObj, parentUrl: string, filter: Object): tRouteList {
  const routeUrls = Object.keys(routesObj).filter((key) => /^\//.test(key));
  const routeList = routeUrls.map(u => {
    const name = u.slice(1);
    return filter[name] ? {
      ...(routesObj[u] instanceof Object && !(routesObj[u] instanceof Array) ? routesObj[u] : {}),
      url: (parentUrl || '') + u,
      name,
      role: filter[name] instanceof Array ? filter[name] : []
    } : null;
  });
  const renderedRoute = routeList.map((route) => {
    if (!route) {
      return null;
    }
    const childRoutes = parseRouteTree(route, route.url, filter[route.name]);
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
