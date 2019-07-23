// @flow

import { USER} from './action';
import defaultAvatarImg from '../../../resources/imgs/image_placeholder.jpg';
import type {tCommonActionType} from '../../common/type';
import type { tAuthUserInfo } from './action';

const lodash = require('lodash');

const defaultUsers = [];

export default function users(
  state: Array<tAuthUserInfo> = defaultUsers,
  action: tCommonActionType & tAuthUserInfo
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
