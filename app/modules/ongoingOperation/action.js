// @flow

export const ONGOING_OPERATION = {
  FETCH_OK: 'ONGOING_OPERATION.FETCH_OK',
  CLEAN: 'ONGOING_OPERATION.CLEAN'
};

export function fetchOngoingOperationOK(data) {
  return {
    type: ONGOING_OPERATION.FETCH_OK,
    data
  };
}

export function cleanOngoingOperation() {
  return {
    type: ONGOING_OPERATION.CLEAN
  };
}
