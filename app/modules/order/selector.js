import { ORDER_STATUS } from './model';

export const currentOrder = order => order?.currentOrder;
export const processingIndex = order => order?.processingIndex;
export const stepStatus = step => step?.status;
export const stepType = step => step?.type;
export const viewingIndex = order => order?.viewingIndex;
export const orderSteps = order => currentOrder(order)?.steps;
export const currentOrderLength = order => orderSteps(order)?.length;
export const processingStep = order => orderSteps(order)?.[processingIndex(order)];
export const viewingStep = order => orderSteps(order)?.[viewingIndex(order)];

export const todoOrders = orderList => orderList?.filter((o) => o.status === ORDER_STATUS.TODO || !(o.status));
export const doneOrders = orderList => orderList?.filter((o) => o.status === ORDER_STATUS.DONE);
export const excpOrders = orderList => orderList?.filter(
  (o) => o.status === ORDER_STATUS.PENDING || o.status === ORDER_STATUS.FAIL || o.status === ORDER_STATUS.CANCEL
);
