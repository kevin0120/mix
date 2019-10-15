// @flow
import { ORDER_STATUS } from './constants';
import type {
  tOrder,
  tOrderState,
  tOrderStepIdx,
} from './interface/typeDef';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { tStepStatus } from '../step/interface/typeDef';
import type { tClsOrder } from './Order';

export const workingOrder = (orderState: tOrderState): ?tClsOrder => orderState?.workingOrder;
export const viewingOrder = (orderState: tOrderState): ?tClsOrder => orderState?.viewingOrder;
export const viewingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.viewingIndex;
export const viewingStep = (orderState: tOrderState): ?IWorkStep =>
  orderSteps(viewingOrder(orderState))?.[viewingIndex(orderState)] || null;

export const workingIndex = (order: ?tClsOrder): tOrderStepIdx => order?.workingIndex || 0;
export const workingStep = (order: ?tClsOrder): ?IWorkStep =>
  orderSteps(order)?.[workingIndex(order)] || null;

export const orderSteps = (order: ?tClsOrder): ?Array<IWorkStep> => order?.steps || null;
export const orderLength = (order: ?tClsOrder): number => orderSteps(order)?.length || 0;

export const todoOrders = (orderList: Array<tClsOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.TODO || (o && !o?.status));

export const doingOrders = (orderList: Array<tClsOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.WIP);

export const doneOrders = (orderList: Array<tClsOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => o?.status === ORDER_STATUS.DONE);

export const exceptOrders = (orderList: Array<tClsOrder>): Array<tOrder> =>
  orderList && orderList.filter((o) => (
  o?.status === ORDER_STATUS.PENDING || o?.status === ORDER_STATUS.CANCEL)
  );

export const stepStatus = (step: ?IWorkStep): ?tStepStatus => step?.status;
// eslint-disable-next-line flowtype/no-weak-types
export const stepData = (step: ?IWorkStep): ?Object => step?.data;
// eslint-disable-next-line flowtype/no-weak-types
export const stepPayload = (step: ?IWorkStep): ?Object => step?.payload;
export const timeCost = (step: ?IWorkStep): ?number => step && step.timeCost();


export const isPending = (order: ?tClsOrder): boolean => order?.status === ORDER_STATUS.PENDING || false;
export const isCancel = (order: ?tClsOrder): boolean => order?.status === ORDER_STATUS.CANCEL || false;
export const doable = (order: ?tClsOrder): boolean =>
  (order?.status === ORDER_STATUS.WIP ||
    order?.status === ORDER_STATUS.TODO ||
    order?.status === ORDER_STATUS.DONE ||
    (order && !order.status)) || false;

export const pendingable = (order: ?tClsOrder): boolean =>
  (order?.status && (
    order?.status !== ORDER_STATUS.PENDING &&
    order?.status !== ORDER_STATUS.CANCEL &&
    order?.status !== ORDER_STATUS.DONE)) || false;

export const cancelable = (order: ?tClsOrder): boolean =>
  (order?.status &&
    (order?.status !== ORDER_STATUS.CANCEL &&
      order?.status !== ORDER_STATUS.DONE)) ||
  false;
