// @flow

import { array } from 'prop-types';

export class Order {
  steps: Array<Step>;
}

export const ORDER_STEP_STATUS = {
  READY: 'ORDER_STEP_STATUS_READY',
  ENTERING: 'ORDER_STEP_STATUS_ENTERING',
  DOING: 'ORDER_STEP_STATUS_DOING',
  LEAVING: 'ORDER_STEP_STATUS_LEAVING',
  FAIL: 'ORDER_STEP_STATUS_FAIL',
  FINISHED: 'ORDER_STEP_STATUS_FINISHED'
};

export const ORDER_STATUS={
  TODO:'TODO',
  WIP:'WIP',
  DONE:'DONE',
  CANCEL:'CANCEL',
  PENDING:'PENDING',
};

export type Step = {
  +name: string,
  info: string,
  status: string,
  +type: string, // check,collect,instruct,enable,...
  payload: {},
  stash: {},
  steps: Array
};


export function hasSubStep(step) {
  return step.steps instanceof Array && step.steps.length > 0;
}

export function HasValidType(step): boolean {

}


