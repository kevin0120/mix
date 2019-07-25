// @flow

export class Order {
  steps: Array<Step>;
}



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
