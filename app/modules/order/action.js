// @flow
import type { tOrder, tOrderStepIdx, tOrderListData, tOrderStatus } from './interface/typeDef';
import type { IOrder } from './interface/IOrder';
import { ORDER_STATUS, ORDER } from './constants';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { tAnyStepStatus, tStepStatus } from '../step/interface/typeDef';

export type tActUpdateState = {
  type: string,
  step: IOrder | IWorkStep,
  status: tOrderStatus
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
  newOrder: (list: Array<tOrder>) => ({
    type: ORDER.NEW,
    list
  }),
  getListSuccess: (list: Array<tOrderListData>) => ({
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
  workOn: (order: IOrder): tActOrderTrigger => ({
    type: ORDER.WORK_ON,
    order
  }),
  finishOrder: (order: IOrder): tActOrderTrigger => ({
    type: ORDER.FINISH,
    order
  }),
  cancelOrder: (order: IOrder): tActUpdateState => ({
    type: ORDER.STEP.STATUS,
    step: order,
    status: ORDER_STATUS.CANCEL
  }),
  pendingOrder: (order: IOrder): tActUpdateState => ({
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
  stepStatus: (step: IWorkStep, status: tAnyStepStatus, msg: string = '') => ({
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
  finishStep: (step: IWorkStep) => ({
    type: ORDER.STEP.FINISH,
    step
  })
};
