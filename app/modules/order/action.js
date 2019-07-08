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
    ENTER: 'ORDER_STEP_ENTER',
    ENTERED: 'ORDER_STEP_ENTERED',
    LEAVE: 'ORDER_STEP_LEAVE',
    FINISH: 'ORDER_STEP_FINISH',
    FAIL: 'ORDER_STEP_FAIL',
    RESET: 'ORDER_STEP_RESET',
    // 步进、步退
    TRY_PUSH: 'ORDER_STEP_TRY_PUSH',
    PUSH: 'ORDER_STEP_PUSH',
    REVOKE: 'ORDER_STEP_REVOKE'
  },
  // JOB:{
  //   START:'ORDER_JOB_START',
  //   SUCCESS:'ORDER_JOB_SUCCESS',
  //   FAIL:'ORDER_JOB_FAIL',
  // },
  CHECK: {
    START: 'ORDER_CHECK_START',
    SUCCESS: 'ORDER_CHECK_SUCCESS',
    FAIL: 'ORDER_CHECK_FAIL'
  },
  COLLECT: {
    START: 'ORDER_COLLECT_START',
    SUCCESS: 'ORDER_COLLECT_SUCCESS',
    FAIL: 'ORDER_COLLECT_FAIL'
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
  finishOrder:()=>({
    type:ORDER.FINISH
  }),
  failOrder:()=>({
    type:ORDER.FAIL
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
  enterStep: () => ({
    type: ORDER.STEP.ENTER
  }),
  enteredStep: () => ({
    type: ORDER.STEP.ENTERED
  }),
  leaveStep: () => ({
    type: ORDER.STEP.LEAVE
  }),
  finishStep: () => ({
    type: ORDER.STEP.FINISH
  }),
  failStep: () => ({
    type: ORDER.STEP.FAIL
  }),
  resetStep: () => ({
    type: ORDER.STEP.RESET
  }),
  pushStep: () => ({
    type: ORDER.STEP.PUSH
  }),
  tryPushStep: () => ({
    type: ORDER.STEP.TRY_PUSH
  })

};

// export const orderStartJob = () => ({
//   type: ORDER.JOB.START
// });
//
// export const orderRevokeJob = () => ({
//   type: ORDER.JOB.REVOKE
// });
