/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */


// @flow

import { call, take } from 'redux-saga/effects';
import { HEALTHZ_CHECK } from '../actions/actionTypes'


function startHealthzCheck() {
  return setInterval(() => {
    console.log('healthz check demo')
  }, 3000);
}

function stopHealthzCheck(timer) {
  clearInterval(timer);
}

export function* healthzCheckFlow() {
  while (yield take(HEALTHZ_CHECK.START)) {
    const timer = yield call(startHealthzCheck);
    yield take(HEALTHZ_CHECK.STOP);
    yield call(stopHealthzCheck, timer)
  }

}
