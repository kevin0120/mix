// @flow
import type { tOrder, tOrderStepIdx, tStepStatus } from './interface/typeDef';
import type { tClsOrder } from './Order';
import type { tClsStep } from '../step/Step';
import { ORDER_STATUS, ORDER } from './constants';

export type tActUpdateState = {
  type: string,
  step: tClsOrder | tClsStep,
  status: tStepStatus
};

export type tActOrderTrigger = {
  type: string,
  order: tOrder
};

export const orderActions = {
  getList: () => ({
    type: ORDER.LIST.GET
  }),
  // eslint-disable-next-line flowtype/no-weak-types
  newOrder: (list: Array<any>) => ({
    type: ORDER.NEW,
    list
  }),
  getListSuccess: (list: Array<tOrder>) => ({
    type: ORDER.LIST.SUCCESS,
    list
  }),
  getListFail: () => ({
    type: ORDER.LIST.FAIL
  }),
  getDetail: (order: tOrder) => ({
    type: ORDER.DETAIL.GET,
    order
  }),
  getDetailSuccess: (order: tOrder) => ({
    type: ORDER.DETAIL.SUCCESS,
    order
  }),
  getDetailFail: () => ({
    type: ORDER.DETAIL.FAIL
  }),
  view: (order: tOrder) => ({
    type: ORDER.VIEW,
    order
  }),
  // order status
  workOn: (order: tClsOrder): tActOrderTrigger => ({
    type: ORDER.WORK_ON,
    order
  }),
  finishOrder: (order: tClsOrder): tActOrderTrigger => ({
    type: ORDER.FINISH,
    order
  }),
  cancelOrder: (order: tClsOrder): tActUpdateState => ({
    type: ORDER.STEP.STATUS,
    step: order,
    status: ORDER_STATUS.CANCEL
  }),
  pendingOrder: (order: tClsOrder): tActUpdateState => ({
    type: ORDER.STEP.STATUS,
    step: order,
    status: ORDER_STATUS.PENDING
  }),
  nextStep: () => ({
    type: ORDER.STEP.NEXT
  }),
  previousStep: () => ({
    type: ORDER.STEP.PREVIOUS
  }),
  jumpToStep: (stepId: tOrderStepIdx) => ({
    type: ORDER.STEP.JUMP_TO,
    stepId
  }),
  stepStatus: (step: tClsStep, status: tStepStatus, msg: string = '') => ({
    type: ORDER.STEP.STATUS,
    step,
    status,
    msg
  }),
  doNextStep: () => ({
    type: ORDER.STEP.DO_NEXT
  }),
  doPreviousStep: () => ({
    type: ORDER.STEP.DO_PREVIOUS
  }),
  updateState: () => ({
    type: ORDER.UPDATE_STATE
  }),
  finishStep: (step: tClsStep) => ({
    type: ORDER.STEP.FINISH,
    step
  })
};
