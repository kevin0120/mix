// @flow

export class Order {
  steps: Array<Step>;
}

export const STEP_STATUS = {
  READY: 'STEP_STATUS_READY',
  ENTERING: 'STEP_STATUS_ENTERING',
  DOING: 'STEP_STATUS_DOING',
  LEAVING: 'STEP_STATUS_LEAVING',
  FAIL: 'STEP_STATUS_FAIL',
  FINISHED: 'STEP_STATUS_FINISHED'
};

export const ORDER_STATUS = {
  TODO: 'TODO',
  WIP: 'WIP',
  DONE: 'DONE',
  CANCEL: 'CANCEL',
  PENDING: 'PENDING',
  FAIL: 'FAIL'
};

export type Step = {
  +name: string,
  info: string,
  status: string,
  +type: string, // check,collect,instruct,enable,...
  payload: {},
  data: {},
  steps: Array
};
