// @flow

import type { Node } from 'react';
import type { IWorkStep } from '../../modules/step/interface/IWorkStep';

export type tStepProps = {|
  bindAction: ((Node => Node) | Node) => void,
  step: IWorkStep,
  isCurrent: boolean,
  bindDescription: ((Node => Node) | Node) => void
|};
