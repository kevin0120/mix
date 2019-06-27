/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

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
