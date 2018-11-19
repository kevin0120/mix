/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { call, take, put } from 'redux-saga/effects';

import { CONNECTION, SYSTEM_INIT } from '../actions/actionTypes';

import { fetchConnectionInfo } from './api/systemInit';
import { initRush } from '../actions/rush';

export function* fetchConnectionFlow(baseUrl, hmiSN, dispatch) {
  const fullUrl = `${baseUrl}/hmi.connections/${hmiSN}`;
  const resp = yield call(fetchConnectionInfo, fullUrl);

  if (resp.status === 200) {
    yield put({ type: CONNECTION.FETCH_OK, data: resp.data });

    // 初始化rush
    yield call(initRush, dispatch, resp.data.masterpc.connection, hmiSN);
  }
}

export function* sysInitFlow() {
  const { baseUrl, hmiSN, dispatch } = yield take(SYSTEM_INIT); // 只获取一次
  try {
    yield call(fetchConnectionFlow, baseUrl, hmiSN, dispatch);
  } catch (e) {
    console.log(e);
  }
}
