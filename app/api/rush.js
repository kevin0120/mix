import { CommonLog } from '../common/utils';

import {isNil} from 'lodash-es';

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

if(!isNil(ipcRenderer)){
  ipcRenderer.on('rush-reply', (event, args, replySN)=>{
    if(messageSNs[replySN]){
      CommonLog.Info('rush-reply', args);
      messageSNs[replySN](args);
      delete messageSNs[replySN];
    }
  });
}


export function rushSendApi(msgType, data, timeout = defaultTimeout) {
  const sn = getSN();
  return new Promise((resolve) => {
    messageSNs[sn] = (args)=>{
      resolve(args);
    };
    ipcRenderer.send('rush-send', {
      data: {
        type: msgType,
        data
      }, timeout, sn
    });
  }).then(resp => resp);
}



