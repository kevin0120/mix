// @flow

import { take, put, call, fork, select, takeEvery } from 'redux-saga/effects';

import type { Saga } from 'redux-saga';

import { push } from 'connected-react-router';

import { getUserInfo } from '../../api/user';

import { loginSuccess, logoutSuccess, USER } from './action';

import type{ tAuthUserInfo } from './action';


import { setNewNotification } from '../notification/action';

const lodash = require('lodash');


type tAuthRespData = {
  +id: number,
  +name: string,
  +uuid: string,
  +image_small: string
};

type tAuthInfo = {
  +user: string | null,
  +password: string | null,
  +method: string
};

type tAuthLogout = {
  +user: string | null
};

function* authorize(user: string, password: string) {
  try {
    const state = yield select();
    const { setting, users, systemSettings } = state;
    const { authEnable } = systemSettings;
    if (authEnable && user === '') {
      // 强制需要认证
      return;
    }
    const u = user !== '' ? user : setting.base.userInfo.uuid;
    const isExisted =
      lodash.some(users, { uuid: u }) || lodash.some(users, { name: u }); // 检测是否已经登录
    if (isExisted) return;

    const response = yield call(
      getUserInfo,
      setting.page.odooConnection.odooUrl.value,
      u,
      password
    );
    const statusCode = response.status;
    if (statusCode === 200) {
      const { id, name, uuid, image_small: avatar }: tAuthRespData = response.data;
      const userInfo: tAuthUserInfo = {
        uid: id,
        name,
        uuid,
        avatar,
        role: 'admin'
      };
      yield put(loginSuccess(userInfo));
      yield put(push('/app'));
    }
  } catch (e) {
    yield put(setNewNotification('error', e));
  }
}

function* logout(action: tAuthLogout): Saga<void> {
  try {
    const { user } = action;
    const state = yield select();
    const { users } = state;
    const deepUsers = lodash.cloneDeep(users);
    const userInfo: tAuthUserInfo = lodash.find(deepUsers, i => i.name === user || i.uuid === user); // 检测是否已经登录
    if (lodash.isUndefined(userInfo)) {
      // 未找到
      return;
    }
    const ret = lodash.remove(state, i => i.name === user || i.uuid === user); // 尝试删除，确认是否要跳转到登录页面
    if (ret.length === 0) {
      // 回到登录页面
      yield put(push('/pages/login'));
    }
    yield put(logoutSuccess(userInfo));
  } catch (e) {
    console.error(e);
  }
}

const loginMethodMap = {
  local: loginLocal,
  online: authorize
};

function* loginLocal(user, password) {
  try {
    const state = yield select();
    const { localUsers } = state.setting.authorization;
    const success = !!localUsers[user] && localUsers[user].password === password;

    if (success) {
      const userInfo: tAuthUserInfo = {
        uid: localUsers[user].uid,
        name: user,
        uuid: localUsers[user].uuid,
        avatar: localUsers[user].avatar,
        role: localUsers[user].role
      };
      yield put(loginSuccess(userInfo));
      yield put(push('/app'));

    }

  } catch (e) {
    console.error(e);
    yield put(setNewNotification('local login error', e));
  }
}

export function* loginFlow(): Saga<void> {
  try {
    while (true) {
      const { user, password, method }: tAuthInfo = yield take(USER.LOGIN.REQUEST);
      yield fork(loginMethodMap[method], user, password);
    }
  } catch (e) {
    console.error(e);
  }
}

export function* logoutFlow(): Saga<void> {
  yield takeEvery(USER.LOGOUT.REQUEST, logout);
}
