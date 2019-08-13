// @flow
import type { tOrder, tOrderStepIdx, tStepStatus } from './model';

export const ORDER = {
  WORK_ON: 'ORDER_WORK_ON',
  VIEW: 'ORDER_VIEW',
  STATUS:'ORDER_STATUS',
  FINISH: 'ORDER_FINISH',
  CANCEL: 'ORDER_CANCEL',
  PENDING: 'ORDER_PENDING',
  LIST: {
    GET: 'ORDER_LIST_GET',
    SUCCESS: 'ORDER_LIST_SUCCESS',
    FAIL: 'ORDER_LIST_FAIL'
  },
  DETAIL: {
    GET: 'ORDER_DETAIL_GET',
    SUCCESS: 'ORDER_DETAIL_SUCCESS',
    FAIL: 'ORDER_DETAIL_FAIL'
  },
  STEP: {
    // 仅移动指针，不修改step状态
    NEXT: 'ORDER_STEP_NEXT',
    PREVIOUS: 'ORDER_STEP_PREVIOUS',
    VIEW_PREVIOUS: 'ORDER_STEP_VIEW_PREVIOUS', // 防抖后事件
    VIEW_NEXT: 'ORDER_STEP_VIEW_NEXT', // 防抖后事件
    JUMP_TO: 'ORDER_STEP_JUMP_TO',
    // 修改step的状态
    STATUS: 'ORDER_STEP_STATUS',
    // 步进、步退
    DO_NEXT: 'ORDER_STEP_DO_NEXT',
    DO_PREVIOUS: 'ORDER_STEP_DO_PREVIOUS',
    // 修改store
    DATA: 'ORDER_STEP_DATA',
    // 记录开始时间
    TIME: 'ORDER_STEP_TIME'
  }
};

export type orderTriggerType = {
  type: string,
  order: tOrder
};

export const orderActions = {
  getList: () => ({
    type: ORDER.LIST.GET
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
  workOn: (order: tOrder): orderTriggerType => ({
    type: ORDER.WORK_ON,
    order
  }),
  finishOrder: () => ({
    type: ORDER.FINISH
  }),
  cancelOrder: () => ({
    type: ORDER.CANCEL
  }),
  pendingOrder: (order: tOrder) => ({
    type: ORDER.PENDING,
    order
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
  stepStatus: (status: tStepStatus) => ({
    type: ORDER.STEP.STATUS,
    status
  }),
  doNextStep: () => ({
    type: ORDER.STEP.DO_NEXT
  }),
  doPreviousStep: () => ({
    type: ORDER.STEP.DO_PREVIOUS
  }),
  stepData: (reducer: ({})=>{}) => ({
    type: ORDER.STEP.DATA,
    reducer
  }),
  stepTime: (idx: tOrderStepIdx, time: Date) => ({
    type: ORDER.STEP.TIME,
    idx,
    time
  })
};
