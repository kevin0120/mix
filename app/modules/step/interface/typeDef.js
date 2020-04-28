// @flow

import { stepTypeKeys, STEP_STATUS } from '../constants';
import type { tWorkableData } from '../../workable/typeDef';

export type tStepInfo = {
  [key: string]: ?string | ?number
};

export type tStepPayload = {
  // eslint-disable-next-line flowtype/no-weak-types
  [key: string]: any
};

export type tStep = {
  ...tWorkableData,
  sequence: number,
  test_type: tStepType, //
  failure_msg: string,
  desc: string,
  image: string,
  skippable: boolean, // 此工步是否可跳过
  undoable: boolean, // 是否可重做
  data: any, // 工步执行过程中生成的数据
  status: tStepStatus,
  payload: tStepPayload, // 工步的数据
  steps: Array<$Shape<tStep>>
};

export type tAnyStatus = any;
export type tRunSubStepCallbacks = {
  onExit?: Function,
  onNext?: Function,
  onPrevious?: Function,
  onAnother?: Function
};


export type tStepType = $Values<typeof stepTypeKeys>;
export type tStepStatus = $Values<typeof STEP_STATUS>;
