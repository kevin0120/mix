// @flow

import type { Node } from 'react';
import type { tClsStep } from '../../modules/step/Step';

export type tStepProps = {|
  bindAction: (Node)=>any,
  step: tClsStep,
  isCurrent: boolean,
  bindDescription: (Node)=>any
|};
