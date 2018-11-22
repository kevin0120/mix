// @flow

import configs from '../shared/config';

import { USER_CONFIGS } from '../actions/actionTypes';

type actionType = {
  +type: string,
  +section: string,
  +newConfigs: object
};

export default function setting(state: object = configs, action: actionType) {
  switch (action.type) {
    case USER_CONFIGS.SAVE: {
      const { section, newConfigs } = action;
      return { ...state, page: { ...state.page, [section]: newConfigs } };
    }
    default:
      return state;
  }
}
