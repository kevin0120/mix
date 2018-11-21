/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { USER } from './actionTypes';

export function doUserAuth(user = null, password = null) {
  return {
    type: USER.LOGIN_REQUEST,
    user,
    password
  };
}

export function userLogOut(user = null) {
  return {
    type: USER.LOGOUT,
    user
  };
}
