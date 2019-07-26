import { ORDER_STATUS } from './model';
import type { tOrder, tOrderState, tOrderStepIdx, tStep } from './model';

export const currentOrder = (orderState: tOrderState): tOrder => orderState?.currentOrder;
export const processingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.processingIndex;
export const viewingIndex = (orderState: tOrderState): tOrderStepIdx => orderState?.viewingIndex;
export const orderSteps = (orderState: tOrderState): Array<tStep> => currentOrder(orderState)?.steps;
export const orderLength = (orderState: tOrderState): number => orderSteps(orderState)?.length;
export const getStep = (orderState: tOrderState, idx: tOrderStepIdx): tStep => orderSteps(orderState)?.[idx];
export const processingStep = (orderState: tOrderState): tStep => orderSteps(orderState)?.[processingIndex(orderState)];
export const viewingStep = (orderState: tOrderState): tStep => orderSteps(orderState)?.[viewingIndex(orderState)];


export const todoOrders = (orderList: Array<tOrder>): Array<tOrder> => orderList?.filter((o) => o.status === ORDER_STATUS.TODO || !(o.status));
export const doneOrders = (orderList: Array<tOrder>): Array<tOrder> => orderList?.filter((o) => o.status === ORDER_STATUS.DONE);
export const exceptOrders = (orderList: Array<tOrder>): Array<tOrder> => orderList?.filter(
  (o) => o.status === ORDER_STATUS.PENDING || o.status === ORDER_STATUS.FAIL || o.status === ORDER_STATUS.CANCEL
);

export const stepStatus = (step: tStep): string => step?.status;
export const stepType = (step: tStep): string => step?.type;
export const stepData = (step: tStep): {} => step?.data;
export const stepPayload = (step: tStep): {} => step?.payload;
export const startTime = (step: tStep): Date => step?.startTime;
export const endTime = (step: tStep): Date => step?.endTime;
