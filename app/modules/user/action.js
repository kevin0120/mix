// @flow

import type { rActionUserType } from './reducer';

export const USER = {
  LOGIN_REQUEST: 'USER.LOGIN_REQUEST',
  LOGIN_SUCCESS: 'USER.LOGIN_SUCCESS',
  LOGOUT: 'USER.LOGOUT',
  LOGOUT_SUCCESS: 'USER.LOGOUT_SUCCESS'
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

export function doUserAuth(user: string = '', password: string = ''): tAuthLogin {
  return {
    type: USER.LOGIN_REQUEST,
    user,
    password
  };
}

export function loginSuccess(data: tAuthUserInfo): rActionUserType {
  return {
    type: USER.LOGIN_REQUEST,
    ...data
  };
}

export function logoutSuccess(data: tAuthUserInfo): rActionUserType {
  return {
    type: USER.LOGOUT_SUCCESS,
    ...data
  };
}

export function userLogOut(user: string = '') {
  return {
    type: USER.LOGOUT,
    user
  };
}
