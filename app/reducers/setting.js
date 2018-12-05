// @flow

import configs from '../shared/config';

import { USER_CONFIGS } from '../actions/actionTypes';

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
      switch (section){
        case 'connections':
          return { ...state, system: { ...state.system, [section]: newConfigs } };
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
    default:
      return state;
  }
}
