// @flow

import { array } from 'prop-types';
import { USER } from './action';

import defaultAvatarImg from '../../../resources/imgs/image_placeholder.jpg';
import { push } from "connected-react-router";

const lodash = require('lodash');

export const defaultUsers = [
  {
    uuid: '11',
    name: 'dummy',
    avatar: defaultAvatarImg,
    uid: 10,
    role: 'admin'
  }
];

export type rActionUserType = {
  +type: string,
  +name: string,
  +avatar: string,
  +uid: number,
  +uuid: string,
  +role: string
};

export default function users(state: array = defaultUsers, action: rActionUserType) {
  switch (action.type) {
    case USER.LOGIN_SUCCESS: {
      const { uid, name, uuid, avatar, role } = action;
      const img =
        lodash.isNil(avatar) || avatar === '' ? defaultAvatarImg : avatar;
      if (state.length === 1 && state[0].name === 'dummy') {
        // 默认用户
        return [{ uid, name, uuid, avatar: img, role }];
      }
      return [...state, { uid, name, uuid, avatar: img, role }];
    }
    case USER.LOGOUT_SUCCESS: {
      const { uid } = action;
      return lodash.remove(state, i => i.name === uid || i.uuid === uid);
    }
    default:
      return state;
  }
}
