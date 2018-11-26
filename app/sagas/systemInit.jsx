/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { call, take, put, select } from 'redux-saga/effects';

import { CONNECTION, SYSTEM_INIT, RUSH, RFID } from '../actions/actionTypes';

import { fetchConnectionInfo } from './api/systemInit';
import {setLedStatusReady} from './io';

// import { initRush } from '../actions/rush';
import { startHealthzCheck } from '../actions/healthCheck';
import { setNewNotification } from '../actions/notification';
import { initIO } from '../actions/ioModbus';
import {initAiis}from '../actions/aiis'

export function* fetchConnectionFlow(baseUrl, hmiSN) {
  const fullUrl = `${baseUrl}/hmi.connections/${hmiSN}`;
  try {
    const resp = yield call(fetchConnectionInfo, fullUrl);

    if (resp.status === 200) {
      const state = yield select();
      yield put({ type: CONNECTION.FETCH_OK, data: resp.data });

      const url = resp.data.masterpc.connection;
      const controllers = resp.data.controllers.map(i => i.serial_no);

      yield put(startHealthzCheck(url, controllers)); // 启动healthzcheck 定时器

      // 初始化rush
      yield put({ type: RUSH.INIT });

      // 初始化io
      if (state.setting.systemSettings.modbusEnable) {
        yield put(initIO());
      }

      // 初始化rfid
      if (state.setting.systemSettings.rfidEnabled) {
        yield put({type: RFID.INIT});
      }

      // 初始化aiis(andon)
      if (state.setting.systemSettings.andonEnable) {
        yield put(initAiis(state.setting.page.odooConnection.aiisUrl.value, hmiSN));
      }
    }

    setLedStatusReady();
  }catch (e) {
    yield put(setNewNotification('error', e.toString()));
  }

}

export function* sysInitFlow() {
  while(true) {
    try {
      const { baseUrl, hmiSN } = yield take(SYSTEM_INIT); // 只获取一次
      yield call(fetchConnectionFlow, baseUrl, hmiSN);
    } catch (e) {
      yield put(setNewNotification('error', e.toString()));
    }
  }
}
