// @flow

import {
  take,
  put,
  call,
  fork,
  select,
  takeEvery,
  all
} from 'redux-saga/effects';

import type { Saga } from 'redux-saga';

import { push } from 'connected-react-router';

import { getUserInfo } from '../../api/user';

import { loginSuccess, logoutSuccess, USER } from './action';

import notifierActions from '../Notifier/action';
import type { tUser } from './interface/typeDef';
import { CommonLog } from '../../common/utils';

const lodash = require('lodash');

type tAuthRespData = {
  +id: number,
  +name: string,
  +uuid: string,
  +image_small: string
};

type tAuthInfo = {
  +name: string | null,
  +password: string | null,
  +method: string
};

type tAuthLogout = {
  +uuid: string | null
};

function* authorize(action) {
  try {
    const { name, password, uuid } = action;
    if (uuid) {
      const state = yield select();
      const { setting } = state;
      const response = yield call(
        getUserInfo,
        setting.page.odooConnection.odooUrl.value,
        uuid
      );
      const statusCode = response.status;
      if (statusCode === 200) {
        const {
          id,
          name: n,
          uuid,
          image_small: avatar
        }: tAuthRespData = response.data;
        const userInfo: tUser = {
          uid: id,
          name: n,
          uuid,
          avatar,
          role: 'admin'
        };
        yield put(loginSuccess(userInfo));
        const newState = yield select();
        if (!/\/app/.test(newState.router.location.pathname)) {
          yield put(push('/app'));
        }
      }
      return;
    }

    if (name && password) {
      const state = yield select();
      const { setting, users, systemSettings } = state;
      const { authEnable } = systemSettings;
      if (authEnable && name === '') {
        // 强制需要认证
        return;
      }
      const u = name !== '' ? name : setting.base.userInfo.uuid;
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
        const {
          id,
          name: n,
          uuid,
          image_small: avatar
        }: tAuthRespData = response.data;
        const userInfo: tUser = {
          uid: id,
          n,
          uuid,
          avatar,
          role: 'admin'
        };
        yield put(loginSuccess(userInfo));
        const newState = yield select();
        if (!/\/app/.test(newState.router.location.pathname)) {
          yield put(push('/app'));
        }
      }
    }
  } catch (e) {
    yield put(notifierActions.enqueueSnackbar('Error', e));
  }
}

function* logout(action: tAuthLogout): Saga<void> {
  try {
    const { uuid } = action;
    const users = yield select(s => s.users);
    const deepUsers = lodash.cloneDeep(users);
    const userInfo: tUser = lodash.find(deepUsers, i => i.uuid === uuid); // 检测是否已经登录
    if (lodash.isUndefined(userInfo)) {
      // 未找到
      return;
    }
    lodash.remove(deepUsers, i => i.uuid === uuid); // 尝试删除，确认是否要跳转到登录页面
    if (deepUsers.length === 0) {
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

function* loginLocal(action) {
  try {
    const { name, password, uuid } = action;
    if (name && password) {
      const state = yield select();
      const { localUsers } = state.setting.authorization;
      const success =
        !!localUsers[name] && localUsers[name].password === password;
      if (success) {
        const userInfo: tUser = {
          uid: localUsers[name].uid,
          name,
          uuid: localUsers[name].uuid,
          avatar: localUsers[name].avatar,
          role: localUsers[name].role
        };
        yield put(loginSuccess(userInfo));
        yield put(push('/app'));
      }
      return;
    }
    if (uuid) {
      const state = yield select();
      const { localUsers } = state.setting.authorization;
      let user = null;
      let n = null;
      Object.keys(localUsers).forEach(k => {
        if (localUsers[k].uuid === uuid) {
          user = localUsers[k];
          n = k;
        }
      });
      if (user) {
        const userInfo: tUser = {
          uid: user.uid,
          name: n,
          uuid: user.uuid,
          avatar: user.avatar,
          role: user.role
        };
        yield put(loginSuccess(userInfo));
        yield put(push('/app'));
      }
    }
  } catch (e) {
    CommonLog.lError(e);
    yield put(notifierActions.enqueueSnackbar('Error', e));
  }
}

export function* loginFlow(): Saga<void> {
  try {
    while (true) {
      const action: tAuthInfo = yield take(USER.LOGIN.REQUEST);
      yield fork(loginMethodMap[action.method], action);
    }
  } catch (e) {
    console.error(e);
  }
}

export function* logoutFlow(): Saga<void> {
  yield takeEvery(USER.LOGOUT.REQUEST, logout);
}

export default function* userRoot(): Saga<void> {
  try {
    yield all([call(loginFlow), call(logoutFlow)]);
  } catch (e) {
    console.error(e);
  }
}
