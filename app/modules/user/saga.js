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

import { isNil, some, cloneDeep, find, isUndefined, remove } from 'lodash-es';
import status from 'http-status';
import { getUserInfo } from '../../api/user';

import {
  loginRequestUuid,
  loginSuccess,
  logoutSuccess,
  USER,
  userNewReader
} from './action';
import notifierActions from '../Notifier/action';
import type {
  tUser,
  tAuthRespData,
  tAuthInfo,
  tAuthLogout
} from './interface/typeDef';
import { CommonLog } from '../../common/utils';
import { bindNewDeviceListener } from '../deviceManager/handlerWSData';
import ClsReader from '../device/reader/ClsReader';
import ioModel from '../io';
import { ioOutputGroups } from '../io/constants';

const DummyUserName = 'DummyUser';

function* authenticate(action) {
  try {
    const { name, password, uuid } = action;
    const state = yield select();
    if (uuid) {
      const { setting } = state;
      const response = yield call(
        getUserInfo,
        setting.page.odooConnection.odooUrl.value,
        uuid
      );
      const statusCode = response.status;
      if (statusCode === status.OK) {
        const {
          id,
          name: respName,
          uuid: respUUID,
          image_small: avatar
        }: tAuthRespData = response.data;
        const userInfo: tUser = {
          uid: id,
          name: respName,
          uuid: respUUID,
          avatar,
          role: 'admin'
        };
        yield put(loginSuccess(userInfo));
        const newState = yield select();
        if (!/\/app/.test(newState.router.location.pathname)) {
          yield put(ioModel.action.setIOOutput({ group: ioOutputGroups.unlock, status: true }));
          yield put(push('/app'));
        }
      }
    } else if (name && password) {
      const { setting, users } = state;
      const { systemSettings } = setting;
      const { authEnable } = systemSettings;
      if (authEnable && name === '') {
        // 强制需要认证
        return;
      }
      const u = name !== '' ? name : setting.base.userInfo.uuid;
      const isExisted = some(users, { uuid: u }) || some(users, { name: u }); // 检测是否已经登录
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
          name: respName,
          uuid: respUUID,
          image_small: avatar
        }: tAuthRespData = response.data;
        const userInfo: tUser = {
          uid: id,
          name: respName,
          uuid: respUUID,
          avatar,
          role: 'admin'
        };
        yield put(loginSuccess(userInfo));
        const newState = yield select();
        if (!/\/app/.test(newState.router.location.pathname)) {
          yield put(ioModel.action.setIOOutput({ group: ioOutputGroups.unlock, status: true }));
          yield put(push('/app'));
        }
      }
    }
  } catch (e) {
    CommonLog.lError(
      `login Workflow User Authentication Error: ${e.toString()}`
    );
    yield put(notifierActions.enqueueSnackbar('Error', e));
  }
}

function* logout(action: tAuthLogout): Saga<void> {
  try {
    const { uuid } = action;
    const users = yield select(s => s.users);
    const deepUsers = cloneDeep(users);
    const userInfo: tUser = find(deepUsers, i => i.uuid === uuid); // 检测是否已经登录
    if (isUndefined(userInfo)) {
      // 未找到
      return;
    }
    remove(deepUsers, i => i.uuid === uuid); // 尝试删除，确认是否要跳转到登录页面
    if (deepUsers.length === 0) {
      // 回到登录页面
      yield put(ioModel.action.setIOOutput({ group: ioOutputGroups.unlock, status: false }));
      yield put(push('/pages/login'));
    }
    yield put(logoutSuccess(userInfo));
  } catch (e) {
    CommonLog.lError(
      `logout Workflow User Authentication Error: ${e.toString()}`
    );
  }
}

const loginMethodMap = {
  local: loginLocal,
  online: authenticate
};

interface ILocalUser {
  [key: string]: tUser;
}

function* loginLocal(action) {
  try {
    const { name, password, uuid } = action;
    if (!name && !password && !uuid) {
      return;
    }
    let userInfo: ?tUser = null;
    const state = yield select();
    if (name && password) {
      const { localUsers } = state.setting.authorization;
      const success =
        !!localUsers[name] && localUsers[name].password === password;
      if (success) {
        userInfo = {
          uid: localUsers[name].uid,
          name,
          uuid: localUsers[name].uuid,
          avatar: localUsers[name].avatar,
          role: localUsers[name].role
        };
        // yield put(loginSuccess(userInfo));
        // yield put(push('/app'));
      }
    } else if (uuid) {
      const {
        localUsers
      }: { localUsers: ILocalUser } = state.setting.authorization;
      let user: ?tUser = null;
      let n: string = DummyUserName;
      Object.keys(localUsers).forEach((k: string) => {
        if (localUsers[k].uuid === uuid) {
          user = localUsers[k];
          n = k;
        }
      });
      if (user) {
        userInfo = {
          uid: user.uid,
          name: n,
          uuid: user.uuid,
          avatar: user.avatar,
          role: user.role
        };
      }
    }
    if (!isNil(userInfo)) {
      yield put(loginSuccess(userInfo));
      yield put(ioModel.action.setIOOutput({ group: ioOutputGroups.unlock, status: true }));
      yield put(push('/app'));
    }
  } catch (e) {
    CommonLog.lError(
      `loginLocal login Workflow login Local User Authentication Error: ${e.toString()}`
    );
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
    CommonLog.lError(
      `login Workflow User Authentication Error: ${e.toString()}`
    );
  }
}

export function* logoutFlow(): Saga<void> {
  yield takeEvery(USER.LOGOUT.REQUEST, logout);
}

function* bindNewReader() {
  try {
    yield call(
      bindNewDeviceListener,
      d => d instanceof ClsReader,
      r => userNewReader(r)
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'bindNewReader' });
  }
}

function onNewReader({ reader }) {
  try {
    reader.addListener(
      () => true,
      i => loginRequestUuid(i.data.data.uid, 'online')
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'onNewReader' });
  }
}

export default function* userRoot(): Saga<void> {
  try {
    yield takeEvery(USER.NEW_READER, onNewReader);
    yield call(bindNewReader);
    yield all([call(loginFlow), call(logoutFlow)]);
  } catch (e) {
    CommonLog.lError(`userRoot User Authentication Error: ${e.toString()}`);
  }
}
