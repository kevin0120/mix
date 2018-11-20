/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { OPERATION } from './actionTypes';

export function switch2Ready() {
  return {
    type: OPERATION.FINISHED
  };
}

export function switch2Doing() {
  return {
    type: OPERATION.STARTED
  };
}

export function operationVerified(data) {
  return {
    type: OPERATION.VERIFIED,
    data: data
  };
}
