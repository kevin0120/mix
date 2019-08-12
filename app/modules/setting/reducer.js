// @flow

import configs from '../../shared/config';

import { USER_CONFIGS } from './action';
import { RFID } from '../external/device/rfid/action';

type actionType = {
  +type: string,
  +section: string,
  +newConfigs: object,
  +uuid: string
};

export default function setting(state: object = configs, action: actionType) {
  switch (action.type) {
    case USER_CONFIGS.SAVE: {
      const { section, newConfigs } = action;
      switch (section) {
        case 'connections':
          return {
            ...state,
            system: { ...state.system, [section]: newConfigs }
          };
        case 'network':
          return state;
        default:
          return { ...state, page: { ...state.page, [section]: newConfigs } };
      }
    }

    case USER_CONFIGS.SET_UUID: {
      const { uuid } = action;
      return {
        ...state,
        base: {
          ...state.base,
          userInfo: {
            ...state.base.userInfo,
            uuid
          }
        }
      };
    }
    case RFID.TOGGLE_STATUS: {
      return {
        ...state,
        systemSettings: {
          ...state.systemSettings,
          rfidEnabled: !state.systemSettings.rfidEnabled
        }
      };
    }
    default:
      return state;
  }
}
