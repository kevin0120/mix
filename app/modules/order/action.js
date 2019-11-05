// @flow
import type { tOrder, tOrderStepIdx, tOrderStatus } from './interface/typeDef';
import type { IOrder } from './interface/IOrder';
import { ORDER_STATUS, ORDER } from './constants';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { tAnyStepStatus } from '../step/interface/typeDef';
import type { IScanner } from '../external/device/scanner/IScanner';

export type tActUpdateState = {
  type: string,
  step: IOrder | IWorkStep,
  status: tOrderStatus
};

export type tActOrderTrigger = {
  type: string,
  order: IOrder
};

// eslint-disable-next-line import/prefer-default-export
export const orderActions = Object.freeze({
  newList: (list: Array<IOrder>) => ({
    type: ORDER.NEW_LIST,
    list
  }),
  getList: () => ({
    type: ORDER.LIST.GET
  }),
  getListSuccess: () => ({
    type: ORDER.LIST.SUCCESS
  }),
  getListFail: () => ({
    type: ORDER.LIST.FAIL
  }),
  getDetail: (order: tOrder) => ({
    type: ORDER.DETAIL.GET,
    order
  }),
  getDetailSuccess: () => ({
    type: ORDER.DETAIL.SUCCESS
  }),
  getDetailFail: () => ({
    type: ORDER.DETAIL.FAIL
  }),
  view: (order: IOrder) => ({
    type: ORDER.VIEW,
    order
  }),
  tryView: (order: tOrder) => ({
    type: ORDER.TRY_VIEW,
    order
  }),
  tryViewCode: (code: string | number) => ({
    type: ORDER.TRY_VIEW,
    code
  }),
  // order status
  tryWorkOn: (order: IOrder): tActOrderTrigger => ({
    type: ORDER.TRY_WORK_ON,
    order
  }),
  tryWorkOnCode: (code: string | number) => ({
    type: ORDER.TRY_WORK_ON,
    code
  }),
  workOn: (order: IOrder): tActOrderTrigger => ({
    type: ORDER.WORK_ON,
    order
  }),
  finishOrder: (order: IOrder): tActOrderTrigger => ({
    type: ORDER.FINISH,
    order
  }),
  orderDidFinish: () => ({
    type: ORDER.DID_FINISH
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
  doPreviousStep: () => ({
    type: ORDER.STEP.DO_PREVIOUS
  }),
  finishStep: (step: IWorkStep) => ({
    type: ORDER.STEP.FINISH,
    step
  }),
  updateState: () => ({
    type: ORDER.UPDATE_STATE
  }),
  newScanner: (scanner: IScanner) => ({
    type: ORDER.NEW_SCANNER,
    scanner
  }),
  reportFinish: (
    code: string,
    trackCode: string,
    workCenterCode: string,
    productCode: string,
    dateComplete: Date,
    operation: {}
  ) => ({
    type: ORDER.REPORT_FINISH,
    code,
    trackCode,
    productCode,
    workCenterCode,
    dateComplete,
    operation
  })
});
