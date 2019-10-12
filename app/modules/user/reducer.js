// @flow

import { remove } from 'lodash-es';
import { USER } from './action';
// import defaultAvatarImg from '../../../resources/imgs/image_placeholder.jpg';
import type { tCommonActionType } from '../../common/type';
import type { tUser } from './interface/typeDef';

const defaultUsers = [];

export default function users(
  state: Array<tUser> = defaultUsers,
  action: tCommonActionType & tUser
) {
  switch (action.type) {
    case USER.LOGIN.SUCCESS: {
      const { uid, name, uuid, avatar, role }: tUser = action;
      const idx = state.findIndex(u => u.uuid === uuid);
      if (idx >= 0) {
        return [...state];
      }
      return [...state, { uid, name, uuid, avatar, role }];
    }
    case USER.LOGOUT.SUCCESS: {
      const { uuid } = action;
      remove(state, i => i.uuid === uuid);
      return [...state];
    }
    default:
      return state;
  }
}
