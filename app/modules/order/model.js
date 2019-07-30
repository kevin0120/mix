// @flow

export type tOrder = {
  steps: Array<tStep>,
  status: string
};

export type tOrderStepIdx = number;

export type tOrderState = {
  currentOrder: tOrder,
  processingIndex: tOrderStepIdx,
  viewingIndex: tOrderStepIdx,
  list: Array<tOrder>
};


export const ORDER_STATUS = {
  TODO: 'TODO',
  WIP: 'WIP',
  DONE: 'DONE',
  CANCEL: 'CANCEL',
  PENDING: 'PENDING',
};

export type tStep = {
  +name: string,
  info: string,
  status: string,
  +type: string, // check,collect,instruct,enable,...
  payload: {},
  data: {},
  steps: Array,
  startTime: Date,
  endTime: Date
};
