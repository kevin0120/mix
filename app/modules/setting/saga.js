// @flow

import { select, takeEvery, put } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { USER_CONFIGS } from './action';
import { CONNECTION } from '../connections/action';
import notifierActions from '../Notifier/action';
import { systemInit } from '../systemInit/action';

const eSetting = require('electron-settings');

function* saveConfiguration(action) {
  const { section, newConfigs } = action.data;
  try {
    const state = yield select();

    const { setting } = state;

    yield put({ type: USER_CONFIGS.SAVE, section, newConfigs });
    yield put(notifierActions.enqueueSnackbar('Info', '配置文件保存成功'));

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
          yield put(
            notifierActions.enqueueSnackbar(
              'Error',
              '保存连接信息,重新初始化失败'
            )
          );
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
  } catch (e) {
    yield put(notifierActions.enqueueSnackbar('Error', '配置文件保存失败'));
  }
}

export default function* watchSettingPreSave(): Saga<void> {
  yield takeEvery(USER_CONFIGS.PRE_SAVE, saveConfiguration);
}
