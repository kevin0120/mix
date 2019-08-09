// @flow
import { call, delay, race } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../common/utils';
import rushSendMessage from '../modules/rush/sendMessage';

const timeout = 10000;

type tIOSn = string;

export function* ioSetApi(sn: tIOSn, index: number, status: string): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_IO_SET',
        data: {
          sn,
          index,
          status
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'ioSetApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi'
    });
  }
}

export function* ioContactApi(sn: tIOSn): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_IO_CONTACT',
        data: {
          sn
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'ioContact timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi'
    });
  }
}

export function* ioStatusApi(sn: tIOSn): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_IO_STATUS',
        data: {
          sn
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'ioStatus timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi'
    });
  }
}
