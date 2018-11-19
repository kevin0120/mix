// @flow

import configs from '../shared/config';

type actionType = {
  +type: string
};

export default function setting(state: object = configs, action: actionType) {
  switch (action.type) {
    default:
      return state;
  }
}
