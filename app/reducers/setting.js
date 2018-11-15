// @flow

import { defaultConfigs } from '../shared/config/defaultConfig';

type actionType = {
  +type: string
};

export default function setting(
  state: object = defaultConfigs,
  action: actionType
) {
  switch (action.type) {
    default:
      return state;
  }
}
