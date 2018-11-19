/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

import { SHUTDOWN_DIAG } from '../actions/actionTypes';

const defaultShutDownDiag = {
  show: false,
  type: '', // 显示的类型
  msg: ''
};

type actionType = {
  +type: string,
  +msg: string,
  +dType: string
};

export default function shutDownDiag(
  state: object = defaultShutDownDiag,
  action: actionType
) {
  switch (action.type) {
    case SHUTDOWN_DIAG.CLOSE: {
      return { ...state, show: false };
    }
    case SHUTDOWN_DIAG.OPEN_WITH_MSG: {
      const { dType, msg } = action.msg;
      return { ...state, show: true, type: dType, msg };
    }
    default:
      return state;
  }
}
