export const OPERATION_VIEWER = {
  LIST_FETCH_START: 'OPERATION_VIEWER_LIST_FETCH_START',
  DETAIL_FETCH_START: 'OPERATION_VIEWER_DETAIL_FETCH_START',
  DETAIL_FETCH_OK: 'OPERATION_VIEWER_DETAIL_FETCH_OK',
  LIST_FETCH_OK: 'OPERATION_VIEWER_LIST_FETCH_OK',
  EDIT_START:'OPERATION_VIEWER_EDIT_START',
  EDIT_END:'OPERATION_VIEWER_EDIT_END',
};
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

export const editOperationEnd = (success) => ({
  type: OPERATION_VIEWER.EDIT_END,
  success
});
