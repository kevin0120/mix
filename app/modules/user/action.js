// @flow

import type { tCommonActionType } from '../../common/type';
import type {
  tUserName,
  tUser,
  tUserLoginAction,
  tUuid
} from './interface/typeDef';

export const USER = {
  LOGIN: {
    REQUEST: 'USER_LOGIN_REQUEST',
    SUCCESS: 'USER_LOGIN_SUCCESS'
  },
  LOGOUT: {
    REQUEST: 'USER_LOGOUT_REQUEST',
    SUCCESS: 'USER_LOGOUT_SUCCESS'
  }
};

type tAuthMethod = 'online' | 'local';

export function loginRequest(
  name: tUserName,
  password: string = '',
  method: tAuthMethod = 'local'
): tUserLoginAction {
  return {
    type: USER.LOGIN.REQUEST,
    name,
    password,
    method
  };
}

export function loginRequestUuid(uuid: string, method: tAuthMethod = 'local') {
  return {
    type: USER.LOGIN.REQUEST,
    uuid,
    method
  };
}

export function loginSuccess(data: tUser): tCommonActionType & tUser {
  return {
    type: USER.LOGIN.SUCCESS,
    ...data
  };
}

export function logoutRequest(uuid: tUuid = '') {
  return {
    type: USER.LOGOUT.REQUEST,
    uuid
  };
}

export function logoutSuccess(data: tUser): tCommonActionType & tUser {
  return {
    type: USER.LOGOUT.SUCCESS,
    ...data
  };
}
