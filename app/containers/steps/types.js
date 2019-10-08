// @flow

import type { Node } from 'react';

export type tStepProps = {
  bindAction: (Node)=>any,
  step: Object,
  isCurrent: boolean,
  bindDescription: (Node)=>any
};
