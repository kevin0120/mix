import type { Node } from 'react';

export type tRouteComponent = Node;

export type tRouteObj = {
  url: tUrl,
  name?: string,
  exact?: boolean,
  component?: tRouteComponent,
  role?: Array<string>
};

export type tUrl = string;

