// @flow
import { Order } from './model';

export const ORDER = {
  TRIGGER: 'ORDER_TRIGGER',
  FINISH: 'ORDER_FINISH',
  FAIL: 'ORDER_FAIL',
  STEP: {
    // 仅移动指针，不修改step状态
    NEXT: 'ORDER_STEP_NEXT',
    PREVIOUS: 'ORDER_STEP_PREVIOUS',
    JUMP_TO: 'ORDER_STEP_JUMP_TO',
    // 修改step的状态
    STATUS: 'ORDER_STEP_STATUS',
    // 步进、步退
    DO_NEXT: 'ORDER_STEP_DO_NEXT',
    DO_PREVIOUS: 'ORDER_STEP_DO_PREVIOUS',
    // 修改store
    DATA:'ORDER_STEP_DATA',
  }
};

export type orderTriggerType = {
  type: string,
  order: Order
};

export const orderActions = {
  trigger: (order: Order): orderTriggerType => ({
    type: ORDER.TRIGGER,
    order
  }),
  finishOrder: () => ({
    type: ORDER.FINISH
  }),
  failOrder: () => ({
    type: ORDER.FAIL
  }),
  nextStep: () => ({
    type: ORDER.STEP.NEXT
  }),
  previousStep: () => ({
    type: ORDER.STEP.PREVIOUS
  }),
  jumpToStep: (stepId) => ({
    type: ORDER.STEP.JUMP_TO,
    stepId
  }),
  stepStatus: () => ({
    type: ORDER.STEP.STATUS
  }),
  doNextStep: () => ({
    type: ORDER.STEP.DO_NEXT
  }),
  doPreviousStep: () => ({
    type: ORDER.STEP.DO_PREVIOUS
  }),
  stepData:(reducer)=>({
    type:ORDER.STEP.DATA,
    reducer
  })
};
