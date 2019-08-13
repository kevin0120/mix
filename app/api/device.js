import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import { CommonLog } from '../common/utils';
import {rushSendApi} from './rush';

export function* deviceStatusApi(): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_DEVICE_STATUS')
  } catch (e) {
    CommonLog.lError(e, {
      at: 'deviceStatusApi',
    });
  }
}