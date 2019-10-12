// @flow
import type { Node } from 'react';
import type { Dispatch } from '../../typeDef';

type tDialogButton = {
  action: Dispatch,
  color: string,
  label: string
};

export type tDialogConfig = {
  buttons: Array<tDialogButton>,
  content: Node,
  title: Node
};

export type tDialogState = {
  open: boolean,
  config: tDialogConfig
};
