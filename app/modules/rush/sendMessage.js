// @flow
import type { Saga } from 'redux-saga';
import { take, fork, call, join } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { getWSClient } from './saga';
import { RUSH } from './action';

const messageSNs = {};

function* listener(sn) {
  try {
    const { payload } = yield take((action) => action.type === RUSH.REPLY && action?.payload?.sn === sn);
    delete messageSNs[sn];
    return payload;
  } catch (e) {
    CommonLog.lError(e);
  }

}

function getSN() {
  let sn = Math.floor(((+new Date()) % 100000000 + Math.random()) * 1000);
  const isSame = k => k === sn;
  while (Object.keys(messageSNs).findIndex(isSame) >= 0) {
    sn = Math.floor(((+new Date()) % 100000000 + Math.random()) * 1000);
  }
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
      ws.sendJson({
        sn,
        payload: data
      }, err => {
        messageSNs[sn] = true;
        if (err && ws) {
          delete messageSNs[sn];
          CommonLog.lError(err);
        }
      });
      return yield join(listenReplyTask);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}



