// @flow
import { tOrder } from './model';

export const ORDER = {
  TRIGGER: 'ORDER_TRIGGER',
  FINISH: 'ORDER_FINISH',
  CANCEL: 'ORDER_CANCEL',
  PENDING: 'ORDER_PENDING',
  SWITCH: 'ORDER_SWITCH',
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
    // 记录开始时间
    START_TIME:'ORDER_STEP_START_TIME',
    END_TIME:'ORDER_STEP_END_TIME'
  }
};

export type orderTriggerType = {
  type: string,
  order: tOrder
};

export const orderActions = {
  trigger: (order: tOrder): orderTriggerType => ({
    type: ORDER.TRIGGER,
    order
  }),
  switchOrder: (order: tOrder): orderTriggerType => ({
    type: ORDER.SWITCH,
    order
  }),
  finishOrder: () => ({
    type: ORDER.FINISH
  }),
  cancelOrder: () => ({
    type: ORDER.CANCEL
  }),
  pendingOrder: () => ({
    type: ORDER.PENDING
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
  stepStatus: (status) => ({
    type: ORDER.STEP.STATUS,
    status
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
  }),
  stepStartTime:(idx,startTime)=>({
    type:ORDER.STEP.START_TIME,
    idx,
    startTime
  }),
  stepEndTime:(idx,endTime)=>({
    type:ORDER.STEP.END_TIME,
    idx,
    endTime
  })
};
