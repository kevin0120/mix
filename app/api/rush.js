import type { Saga } from 'redux-saga';
import { call, delay, race } from 'redux-saga/effects';
import rushSendMessage from '../modules/rush/sendMessage';
import { CommonLog } from '../common/utils';

const defaultTimeout = 10000;

export function* rushSendApi(msgType, data, timeout = defaultTimeout): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: msgType,
        data
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: `rushSendApi timeout ${msgType}`
      };
    }
    const { data: ret } = resp;
    return ret;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'rushSendApi', type: msgType, data
    });
  }
}
