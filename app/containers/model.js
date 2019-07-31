import PropTypes from "prop-types";

export type tRouteComponent = PropTypes.element;

export type tRouteObj = {
  url: tUrl,
  name?: string,
  exact?: boolean,
  component?: tRouteComponent,
  role?: Array<string>
};

export type tUrl = string;

