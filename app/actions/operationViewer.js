import { OPERATION_VIEWER } from './actionTypes';

export const fetchOperationListStart = () => ({
  type: OPERATION_VIEWER.LIST_FETCH_START
});

export const fetchOperationDetailStart = operationID => ({
  type: OPERATION_VIEWER.DETAIL_FETCH_START,
  operationID
});

export const editOperation = (operationID, points, img) => ({
  type: OPERATION_VIEWER.EDIT_START,
  operationID,
  points,
  img
});
