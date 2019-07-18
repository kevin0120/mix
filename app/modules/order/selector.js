export const currentOrder = order => order?.currentOrder;
export const processingIndex = order => order?.processingIndex;
export const stepStatus = step => step?.status;
export const stepType =step=>step?.type;
export const viewingIndex = order => order?.viewingIndex;
export const orderSteps = order => currentOrder(order)?.steps;
export const currentOrderLength = order => orderSteps(order)?.length;
export const processingStep = order => orderSteps(order)?.[processingIndex(order)];
export const viewingStep = order => orderSteps(order)?.[viewingIndex(order)];

