/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import configs from '../shared/config';
import WORK_MODE from '../actions/actionTypes';

const defaultWorkMode = {
  workMode: configs.operationSettings.workMode,
  controllerMode: configs.operationSettings.controllerMode
};

type actionType = {
  +type: string,
  +mode: string
};

export default function workMode(
  state: object = defaultWorkMode,
  action: actionType
) {
  switch (action.type) {
    case WORK_MODE.SWITCH_CM:
      return { ...state, controlMode: action.mode };
    case WORK_MODE.SWITCH_WM:
      return { ...state, workMode: action.mode };
    default:
      return state;
  }
}
