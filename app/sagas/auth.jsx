// @flow

import { take, put, call, fork, select, takeEvery } from 'redux-saga/effects';
import { push, goBack, go } from 'connected-react-router';

import { getUserInfo } from './api/user';

import { USER } from '../actions/actionTypes';

import { setNewNotification } from '../actions/notification';

const lodash = require('lodash');

function* authorize(user, password) {
  const state = yield select();
  const { setting, users } = state;
  const u = user !== null ? user : setting.base.userInfo.uuid;
  const isExisted =
    lodash.some(users, { uuid: u }) || lodash.some(users, { name: u }); // 检测是否已经登录
  if (isExisted) return;
  try {
    const response = yield call(
      getUserInfo,
      setting.page.odooConnection.odooUrl.value,
      u,
      password
    );
    const statusCode = response.status;
    if (statusCode === 200) {
      const { id, name, uuid, image_small } = response.data;
      yield put({
        type: USER.LOGIN_SUCCESS,
        uid: id,
        name,
        uuid,
        avatar: image_small
      });
      yield put(push('/welcome'));
    }
  } catch (e) {
    yield put(setNewNotification('error', e));
  }
}

function* logout(action) {
  const { user } = action;
  const state = yield select();
  const { users } = state;
  let deepUsers = lodash.cloneDeep(users);
  lodash.remove(deepUsers, i => i.name === user || i.uuid === user); // 检测是否已经登录
  if (deepUsers.length === 0) {
    yield put(push('/pages/login'));
  }
  yield put({ type: USER.LOGOUT_SUCCESS, data: deepUsers });
}

export function* loginFlow() {
  while (true) {
    const { user, password } = yield take(USER.LOGIN_REQUEST);
    // fork return a Task object
    yield fork(authorize, user, password);
  }
}

export function* logoutFlow() {
  yield takeEvery(USER.LOGOUT, logout);
}
