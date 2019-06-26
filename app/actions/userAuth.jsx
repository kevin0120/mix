// @flow

import { USER } from './actionTypes';
import type { rActionUserType } from '../reducers/user';

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
