import { ORDER_STATUS } from './model';

export const currentOrder = orderState => orderState?.currentOrder;
export const processingIndex = orderState => orderState?.processingIndex;
export const viewingIndex = orderState => orderState?.viewingIndex;
export const orderSteps = orderState => currentOrder(orderState)?.steps;
export const orderLength = orderState => orderSteps(orderState)?.length;
export const processingStep = orderState => orderSteps(orderState)?.[processingIndex(orderState)];
export const viewingStep = orderState => orderSteps(orderState)?.[viewingIndex(orderState)];


export const todoOrders = orderList => orderList?.filter((o) => o.status === ORDER_STATUS.TODO || !(o.status));
export const doneOrders = orderList => orderList?.filter((o) => o.status === ORDER_STATUS.DONE);
export const exceptOrders = orderList => orderList?.filter(
  (o) => o.status === ORDER_STATUS.PENDING || o.status === ORDER_STATUS.FAIL || o.status === ORDER_STATUS.CANCEL
);

export const stepStatus = step => step?.status;
export const stepType = step => step?.type;
export const stepData = step => step?.data;
