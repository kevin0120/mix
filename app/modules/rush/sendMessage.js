// @flow
import type { Saga } from 'redux-saga';
import { take, fork, join } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { getWSClient } from './client';
import { RUSH } from './action';

const messageSNs = {};

function* listener(sn) {
  try {
    const { payload } = yield take((action) => action.type === RUSH.DATA && action?.payload?.sn === sn);
    delete messageSNs[sn];
    return payload;
  } catch (e) {
    CommonLog.lError(e);
  }

}

function getSN() {
  let sn;
  const isSame = k => k === sn;
  do {
    sn = Math.floor(((+new Date()) % 100000000 + Math.random()) * 1000);
  } while (Object.keys(messageSNs).findIndex(isSame) >= 0);
  return sn;
}

// eslint-disable-next-line flowtype/no-weak-types
export default function* rushSendMessage(data: Object): Saga<void> {
  try {
    const ws = getWSClient();
    if (ws) {
      const sn = getSN();
      // const sn = 1;
      const listenReplyTask = yield fork(listener, sn);
      const msg = {
        sn,
        ...data
      };
      CommonLog.Info('call rush sendMessage', {
        msg: JSON.stringify(msg)
      });
      ws.sendJson(msg, (err) => {
        messageSNs[sn] = true;
        if (err && ws) {
          CommonLog.lError(err);
          delete messageSNs[sn];
        }
      });
      return yield join(listenReplyTask);
    }
    CommonLog.Warn('cannot send message to rush now, rush is not connected');

  } catch (e) {
    CommonLog.lError(e);
  } finally {
    // console.log('rushSendMessage finished');
  }
}



