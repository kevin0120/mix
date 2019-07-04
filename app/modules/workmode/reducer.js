// @flow

import configs from '../../shared/config';
import { WORK_MODE } from './action';

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
