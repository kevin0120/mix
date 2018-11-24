/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { select, takeEvery, put, call } from 'redux-saga/effects';
import { USER_CONFIGS } from '../actions/actionTypes';
import { setNewNotification } from '../actions/notification';
import { initIO } from '../actions/ioModbus';

const eSetting = require('electron-settings');

const getSetting = state => state.setting;

function* saveConfiguration(action) {
  const { section, newConfigs } = action.data;

  const setting = yield select(getSetting);

  try {
    yield put({ type: USER_CONFIGS.SAVE, section, newConfigs });
    eSetting.setAll({ ...setting, [section]: newConfigs });
    yield put(setNewNotification('success', '配置文件保存成功'));
    yield put(initIO());
  } catch (e) {
    yield put(setNewNotification('error', '配置文件保存失败'));
  }

  switch (section) {
    case 'odooConnection': {
      break;
    }
    default:
      break;
  }
}

export function* watchSettingPreSave() {
  yield takeEvery(USER_CONFIGS.PRE_SAVE, saveConfiguration);
}
