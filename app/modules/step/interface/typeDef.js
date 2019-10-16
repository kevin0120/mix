// @flow

import { stepTypeKeys, STEP_STATUS } from '../constants';

export type tStepInfo = {
  [key: string]: ?string | ?number
};

export type tStepPayload = ?{
  // eslint-disable-next-line flowtype/no-weak-types
  [key: string]: any
};

export type tStep = {
  +id: number,
  +name: string,
  +desc: string,
  info: tStepInfo,
  status: tStepStatus,
  +type: tStepType, //
  payload: tStepPayload, // 工步的数据
  data: tStepPayload, // 工步执行过程中生成的数据
  steps: Array<tStep>,
  times: Array<Date>,
  skippable: boolean, // 此工步是否可跳过
  undoable: boolean, // 是否可重做
  description: string
};

export type tStepDataReducer = Function;
export type tAnyStepStatus = string;
export type tRunSubStepCallbacks = {
  onExit?: Function,
  onNext?: Function,
  onPrevious?: Function
};



export type tStepType = $Values<typeof stepTypeKeys>;
export type tStepStatus = $Values<typeof STEP_STATUS>;