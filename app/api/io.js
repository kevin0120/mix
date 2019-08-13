// @flow
import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';

type tIOSn = string;

export function* ioSetApi(sn: tIOSn, index: number, status: string): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_IO_SET', {
      sn,
      index,
      status
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioSetApi'
    });
  }
}

export function* ioContactApi(sn: tIOSn): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_IO_CONTACT', {
      sn
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioContact'
    });
  }
}

export function* ioStatusApi(sn: tIOSn): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_IO_STATUS', {
      sn
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioStatusApi'
    });
  }
}
