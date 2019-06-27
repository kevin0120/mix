/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

export const TOOLS = {
  ENABLE: 'TOOLS.ENABLE',
  DISABLE: 'TOOLS.DISABLE',
  STATUS_CHANGE:'TOOLS.STATUS.CHANGE'
};

export function toolEnable(reason) {
  return {
    type: TOOLS.ENABLE,
    enable: true,
    reason
  };
}

export function toolDisable(reason) {
  return {
    type: TOOLS.ENABLE,
    enable: false,
    reason
  };
}

export function toolStatusChange(toolSN,status,reason){
  return{
    type:TOOLS.STATUS_CHANGE,
    toolSN,
    status,
    reason
  }
}
