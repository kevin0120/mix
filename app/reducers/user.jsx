/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { USER } from '../actions/actionTypes';

import defaultAvatarImg from '../../resources/imgs/image_placeholder.jpg';

const lodash = require('lodash');

export const defaultUsers = [
  {
    uuid: '11',
    name: 'dummy',
    avatar: defaultAvatarImg,
    uid: 10
  }
];

// type actionType = {
//   +type: string,
//   +name: string,
//   +avatar: string,
//   +uid: number,
//   +uuid: string
// };

export default function users(state: Array = defaultUsers, action) {
  switch (action.type) {
    case USER.LOGIN_SUCCESS: {
      const { uid, name, uuid, avatar } = action;
      const img =
        lodash.isNil(avatar) || avatar === '' ? defaultAvatarImg : avatar;
      if (state.length === 1 && state[0].name === 'dummy') {
        return [{ uid, name, uuid, avatar: img }];
      }
      return [...state, { uid, name, uuid, avatar: img }];
    }
    case USER.LOGOUT_SUCCESS: {
      if (action.data.length === 0) {
        return defaultUsers;
      }
      return action.data;
    }
    default:
      return state;
  }
}
