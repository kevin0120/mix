// @flow
import { isNil } from 'lodash-es';
import { CommonLog } from '../common/utils';
import type { tRushData } from '../modules/rush/type';


const { ipcRenderer } = require('electron');

const messageSNs = {};

function getSN() {
  let sn;
  const isSame = k => k === sn;
  do {
    sn = Math.floor(((+new Date()) % 100000000 + Math.random()) * 1000);
  } while (Object.keys(messageSNs).findIndex(isSame) >= 0);
  return sn;
}

const defaultTimeout = 10000; // 默认timeout 10s

if (!isNil(ipcRenderer)) {
  ipcRenderer.on('rush-reply', (event, args, replySN) => {
    CommonLog.Info('rush-reply', args);
    if (messageSNs[replySN]) {
      messageSNs[replySN](args);
      delete messageSNs[replySN];
    }
  });
}


// eslint-disable-next-line flowtype/no-weak-types
export function rushSendApi(msgType: string, data: any, timeout: number = defaultTimeout): Promise<tRushData<any, any>> {
  const sn = getSN();
  return new Promise((resolve, reject) => {
    // eslint-disable-next-line flowtype/no-weak-types
    messageSNs[sn] = (resp: tRushData<any, any>) => {
      if (!resp) {
        reject(new Error(`error (sn:${sn}):ws received an empty response`));
        return;
      }
      if (!resp.data) {
        reject(new Error(`error (sn:${sn}):reply has no data`));
        return;
      }
      if (!isNil(resp.data.result) && resp.data.result < 0) {
        reject(new Error(`error ${resp.data.result}(sn:${sn}):${resp.data.msg}`));
        return;
      }
      resolve(resp);
    };
    ipcRenderer.send('rush-send', {
      data: {
        type: msgType,
        data
      }, timeout, sn
    });
  }).then(resp => resp);
}



