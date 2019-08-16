// @flow
import { ORDER_STATUS } from './model';
import type { tOrder, tOrderState, tOrderStepIdx, tStep, tStepStatus, tStepType } from './model';
import Step from '../step/Step';
import Order from './Order';

export const workingOrder = (orderState: tOrderState): ?Order => orderState?.workingOrder;
export const viewingOrder = (orderState: tOrderState): ?Order => orderState?.viewingOrder;
export const viewingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.viewingIndex;
export const viewingStep = (orderState: tOrderState): ?Step =>
  orderSteps(viewingOrder(orderState))?.[viewingIndex(orderState)] || null;

export const workingIndex = (order: ?Order): tOrderStepIdx => order?.workingIndex || 0;
export const workingStep = (order: ?Order): ?tStep =>
  orderSteps(order)?.[workingIndex(order)] || null;

export const orderSteps = (order: ?Order): ?Array<Step> => order?.steps || null;
export const orderLength = (order: ?Order): number => orderSteps(order)?.length || 0;

export const todoOrders = (orderList: Array<Order>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.TODO || (o && !o?.status));

export const doingOrders = (orderList: Array<Order>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.WIP);

export const doneOrders = (orderList: Array<Order>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.DONE);

export const exceptOrders = (orderList: Array<Order>): Array<tOrder> =>
  orderList && orderList.filter((o) => (
  o?.status === ORDER_STATUS.PENDING || o?.status === ORDER_STATUS.CANCEL)
  );


export const stepStatus = (step: ?tStep): ?tStepStatus => step?.status;
export const stepType = (step: ?tStep): ?tStepType => step?.type;
export const stepData = (step: ?tStep): ?Object => step?.data;
export const stepPayload = (step: ?tStep): ?Object => step?.payload;
export const times = (step: ?Step): ?Array<Date> => step?.times;
export const timeCost = (step: ?Step): ?Array<Date> => step && step.timeCost();


export const isPending = (order: ?Order): boolean => order?.status === ORDER_STATUS.PENDING || false;
export const isCancel = (order: ?Order): boolean => order?.status === ORDER_STATUS.CANCEL || false;
export const doable = (order: ?Order): boolean =>
  (order?.status === ORDER_STATUS.WIP ||
    order?.status === ORDER_STATUS.TODO ||
    (order && !order.status)) || false;

export const pendingable = (order: ?Order): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.PENDING &&
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE)) || false;

export const cancelable = (order: ?tOrder): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE
  )) || false;
