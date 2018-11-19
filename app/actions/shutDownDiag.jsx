/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { SHUTDOWN_DIAG } from './actionTypes';

export function openShutdown(dType, data = null) {
  return {
    type: SHUTDOWN_DIAG.OPEN,
    dType,
    data
  };
}

export function closeShutDownDiag(dType) {
  return {
    type: SHUTDOWN_DIAG.CLOSE_START,
    dType
  };
}

export function confirmShutDownDiag(dType) {
  return {
    type: SHUTDOWN_DIAG.CONFIRM,
    dType
  };
}
