/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { select, takeEvery, put, call } from 'redux-saga/effects';
import { USER_CONFIGS, CONNECTION, NETWORK } from '../actions/actionTypes';
import { setNewNotification } from '../actions/notification';
import { systemInit } from '../actions/sysInit';

const eSetting = require('electron-settings');

const getSetting = state => state.setting;

function* saveConfiguration(action) {
  const { section, newConfigs } = action.data;
  try {
    const state = yield select();

    const { setting } = state;

    yield put({ type: USER_CONFIGS.SAVE, section, newConfigs });
    yield put(setNewNotification('success', '配置文件保存成功'));

    switch (section) {
      case 'connections': {
        eSetting.setAll({
          ...setting,
          system: { ...setting.system, [section]: newConfigs }
        });
        try {
          yield put({ type: CONNECTION.MANUAL_MODIFICATION, data: newConfigs });
          yield put(systemInit(section));
        } catch (e) {
          yield put(setNewNotification('error', '保存连接信息,重新初始化失败'));
        }

        break;
      }
      case 'odooConnection': {
        eSetting.setAll({
          ...setting,
          system: { ...setting.system, [section]: newConfigs }
        });
        break;
      }
      case 'network':
        // yield put({
        //   type:NETWORK.SET,
        //   config:newConfigs
        // });
        eSetting.setAll({
          ...setting,
          page: { ...setting.page, [section]: newConfigs }
        });
        break;
      default:
        eSetting.setAll({
          ...setting,
          page: { ...setting.page, [section]: newConfigs }
        });
        break;
    }
  }catch (e) {
    yield put(setNewNotification('error', '配置文件保存失败'));
  }
}

export function* watchSettingPreSave() {
  yield takeEvery(USER_CONFIGS.PRE_SAVE, saveConfiguration);
}
