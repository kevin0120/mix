// @flow
import { ORDER_STATUS } from './model';
import type { tOrder, tOrderState, tOrderStepIdx, tStep, tStepArray, tStepStatus, tStepType } from './model';

export const workingOrder = (orderState: tOrderState): ?tOrder => orderState?.workingOrder;
export const viewingOrder = (orderState: tOrderState): ?tOrder => orderState?.viewingOrder;
export const workingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.workingIndex;
export const viewingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.viewingIndex;
export const workingStep = (orderState: tOrderState): ?tStep =>
  orderSteps(workingOrder(orderState))?.[workingIndex(orderState)] || null;
export const viewingStep = (orderState: tOrderState): ?tStep =>
  orderSteps(viewingOrder(orderState))?.[viewingIndex(orderState)] || null;


export const orderSteps = (order: ?tOrder): ?tStepArray => order?.steps || null;
export const orderLength = (order: ?tOrder): number => orderSteps(order)?.length || 0;
export const getStep = (order: ?tOrder, idx: tOrderStepIdx): ?tStep => orderSteps(order)?.[idx] || null;

export const todoOrders = (orderList: Array<tOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.TODO || (o && !o?.status));

export const doingOrders = (orderList: Array<tOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.WIP);

export const doneOrders = (orderList: Array<tOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.DONE);

export const exceptOrders = (orderList: Array<tOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => (
  o?.status === ORDER_STATUS.PENDING || o?.status === ORDER_STATUS.CANCEL)
  );


export const stepStatus = (step: ?tStep): ?tStepStatus => step?.status;
export const stepType = (step: ?tStep): ?tStepType => step?.type;
export const stepData = (step: ?tStep): ?Object => step?.data;
export const stepPayload = (step: ?tStep): ?Object => step?.payload;
export const startTime = (step: ?tStep): ?Date => step?.startTime;
export const endTime = (step: ?tStep): ?Date => step?.endTime;

export const isPending = (order: ?tOrder): boolean => order?.status === ORDER_STATUS.PENDING || false;
export const isCancel = (order: ?tOrder): boolean => order?.status === ORDER_STATUS.CANCEL || false;
export const doable = (order: ?tOrder): boolean =>
  (order?.status === ORDER_STATUS.WIP ||
    order?.status === ORDER_STATUS.TODO ||
    (order && !order.status)) || false;

export const pendingable = (order: ?tOrder): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.PENDING &&
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE)) || false;

export const cancelable = (order: ?tOrder): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE
  )) || false;
