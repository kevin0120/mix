// @flow

import type { rActionUserType } from './reducer';

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

type tAuthLogin = {
  +type: string,
  +user: string,
  +password: string
};

export type tAuthUserInfo = {
  +name: string,
  +avatar: string,
  +uid: number,
  +uuid: string,
  +role: string
};

export function loginRequest(user: string = '', password: string = '', method: string = 'local'): tAuthLogin {
  return {
    type: USER.LOGIN.REQUEST,
    user,
    password,
    method
  };
}

export function loginSuccess(data: tAuthUserInfo): rActionUserType {
  return {
    type: USER.LOGIN.SUCCESS,
    ...data
  };
}

export function logoutRequest(user: string = '') {
  return {
    type: USER.LOGOUT.REQUEST,
    user
  };
}

export function logoutSuccess(data: tAuthUserInfo): rActionUserType {
  return {
    type: USER.LOGOUT.SUCCESS,
    ...data
  };
}

