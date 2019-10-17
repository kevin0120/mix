// @flow
import { ORDER_STATUS } from './constants';
import type {
  tOrder,
  tOrderState,
  tOrderStepIdx,
} from './interface/typeDef';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { tAnyStepStatus, tStepStatus } from '../step/interface/typeDef';
import type { IOrder } from './interface/IOrder';

export const workingOrder = (orderState: tOrderState): ?IOrder => orderState?.workingOrder;
export const viewingOrder = (orderState: tOrderState): ?IOrder => orderState?.viewingOrder;
export const viewingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.viewingIndex;
export const viewingStep = (orderState: tOrderState): ?IWorkStep =>
  orderSteps(viewingOrder(orderState))?.[viewingIndex(orderState)] || null;

export const workingIndex = (order: ?IOrder): tOrderStepIdx => order?.workingIndex || 0;
export const workingStep = (order: ?IOrder): ?IWorkStep =>
  orderSteps(order)?.[workingIndex(order)] || null;

export const orderSteps = (order: ?IOrder): ?Array<IWorkStep> => order?.steps || null;
export const orderLength = (order: ?IOrder): number => orderSteps(order)?.length || 0;

export const todoOrders = (orderList: Array<IOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.TODO || (o && !o?.status));

export const doingOrders = (orderList: Array<IOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.WIP);

export const doneOrders = (orderList: Array<IOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.DONE);

export const exceptOrders = (orderList: Array<IOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => (
  o?.status === ORDER_STATUS.PENDING || o?.status === ORDER_STATUS.CANCEL)
  );

export const stepStatus = (step: ?IWorkStep): ?tAnyStepStatus =>step?.status;
// eslint-disable-next-line flowtype/no-weak-types
export const stepData = (step: ?IWorkStep): ?Object => step?.data;
// eslint-disable-next-line flowtype/no-weak-types
export const stepPayload = (step: ?IWorkStep): ?Object => step?.payload;
export const timeCost = (step: ?IWorkStep): ?number => step && step.timeCost();


export const isPending = (order: ?IOrder): boolean => order?.status === ORDER_STATUS.PENDING || false;
export const isCancel = (order: ?IOrder): boolean => order?.status === ORDER_STATUS.CANCEL || false;
export const doable = (order: ?IOrder): boolean =>
  (order?.status === ORDER_STATUS.WIP ||
    order?.status === ORDER_STATUS.TODO ||
    order?.status === ORDER_STATUS.DONE ||
    (order && !order.status)) || false;

export const pendingable = (order: ?IOrder): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.PENDING &&
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE)) || false;

export const cancelable = (order: ?IOrder): boolean =>
  (order?.status &&
    (order?.status !== ORDER_STATUS.CANCEL &&
      order?.status !== ORDER_STATUS.DONE)) ||
  false;
