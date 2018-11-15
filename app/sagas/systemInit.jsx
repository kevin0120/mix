/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { call, take } from 'redux-saga/effects'

import { SYSTEM_INIT } from "../actions/actionTypes";

import { fetchConnectionInfo } from "./api/systemInit";


export function* fetchConnectionFlow(baseUrl, hmiSN) {
  const fullUrl = `${baseUrl}/hmi.connections/${hmiSN}`;
  const info = yield call(fetchConnectionInfo, fullUrl);
  console.log(info);
}

export function* sysInitFlow() {
    const{ baseUrl, hmiSN } = yield take(SYSTEM_INIT); // 只获取一次
  try {
    yield call(fetchConnectionFlow, baseUrl, hmiSN);
  }catch (e) {
    console.log(e);
  }
}
