// @flow

import { USER } from './action';
import defaultAvatarImg from '../../../resources/imgs/image_placeholder.jpg';

const lodash = require('lodash');

const defaultUsers = [];

export type rActionUserType = {
  +type: string,
  +name: string,
  +avatar: string,
  +uid: number,
  +uuid: string,
  +role: string
};

export default function users(
  state: array = defaultUsers,
  action: rActionUserType
) {
  switch (action.type) {
    case USER.LOGIN.SUCCESS: {
      const { uid, name, uuid, avatar, role } = action;
      const img =
        lodash.isNil(avatar) || avatar === '' ? defaultAvatarImg : avatar;
      const idx = state.findIndex(u => u.uid === uid);
      if (idx >= 0) {
        return [...state];
      }
      return [...state, { uid, name, uuid, avatar: img, role }];
    }
    case USER.LOGOUT.SUCCESS: {
      const { uid } = action;
      return lodash.remove(state, i => i.uid === uid);
    }
    default:
      return state;
  }
}
