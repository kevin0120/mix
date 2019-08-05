// @flow

import type {tCommonActionType, tDeviceNewData} from "../../common/type"

export const CONTROLLER = {
  CONTROLLER_SOCKET_DATA: 'CONTROLLER_SOCKET_DATA'
};

export function ControllerNewData(socketInfo: string): tCommonActionType & tDeviceNewData {
  return {
    type: CONTROLLER.CONTROLLER_SOCKET_DATA,
    data: socketInfo,
  };
}
