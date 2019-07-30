// @flow

export type tOrder = {
  steps: tStepArray,
  status: tOrderStatus,
  name: string,
  info: string
};

export type tOrderStepIdx = number;

export type tOrderState = {
  viewingOrder: tOrder | null,
  workingOrder: tOrder | null,
  workingIndex: tOrderStepIdx,
  viewingIndex: tOrderStepIdx,
  list: Array<tOrder>
};

export type tOrderStatus = $Keys<typeof ORDER_STATUS>;

export const ORDER_STATUS = {
  TODO: 'TODO',
  WIP: 'WIP',
  DONE: 'DONE',
  CANCEL: 'CANCEL',
  PENDING: 'PENDING'
};

export type tStepArray = Array<tStep>;

export type tStep = {
  +name: string,
  info: string,
  status: tStepStatus,
  +type: tStepType, // check,collect,instruct,enable,...
  payload: {},
  data: {},
  steps: tStepArray,
  startTime: Date,
  endTime: Date,
  skippable: boolean,
  undoable: boolean,
  description: string
};

export type tStepType = string;


export type tStepStatus = string;
