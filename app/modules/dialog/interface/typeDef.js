// @flow
import PropTypes from 'prop-types';
import type { Dispatch } from '../../typeDef';

type tDialogButton = {
  action: Dispatch,
  color: string,
  label: string
};

export type tDialogConfig = {
  buttons: Array<tDialogButton>,
  content: PropTypes.Element,
  title: PropTypes.Element
};

export type tDialogState = {
  open: boolean,
  config: tDialogConfig
};
