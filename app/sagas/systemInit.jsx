/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { call, take, put, select } from 'redux-saga/effects';

import { CONNECTION, SYSTEM_INIT, RUSH } from '../actions/actionTypes';

import { fetchConnectionInfo } from './api/systemInit';
import { initRush } from '../actions/rush';
import { initIOModbus, setLedStatusReady } from '../actions/ioModbus';
import { startHealthzCheck } from '../actions/healthCheck';
import { setNewNotification } from '../actions/notification';

export function* fetchConnectionFlow(baseUrl, hmiSN, dispatch, getState) {
  const fullUrl = `${baseUrl}/hmi.connections/${hmiSN}`;
  const resp = yield call(fetchConnectionInfo, fullUrl);

  if (resp.status === 200) {
    yield put({ type: CONNECTION.FETCH_OK, data: resp.data });

    const url = resp.data.masterpc.connection;
    const controllers = resp.data.controllers.map(i => i.serial_no);

    yield put(startHealthzCheck(url, controllers)); // 启动healthzcheck 定时器

    // 初始化rush
    // yield call(initRush, dispatch, resp.data.masterpc.connection, hmiSN);
    yield put({ type: RUSH.INIT });

    // 初始化io
    yield call(initIOModbus, dispatch, getState);

    // 初始化aiis
  }

  setLedStatusReady();
}

export function* sysInitFlow() {
  const { baseUrl, hmiSN, dispatch, getState } = yield take(SYSTEM_INIT); // 只获取一次
  try {
    yield call(fetchConnectionFlow, baseUrl, hmiSN, dispatch, getState);
  } catch (e) {
    yield put(setNewNotification('error', e.toString()));
  }
}
