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

const defaultTimeout = 10000;

export function rushSendApi(msgType, data, timeout = defaultTimeout) {
  const sn = getSN();
  messageSNs[sn] = true;
  return new Promise((resolve) => {
    ipcRenderer.send('rush-send', {
      data: {
        type: msgType,
        data
      }, timeout, sn
    });
    const reply = (event, args, replySN) => {
      if (replySN === sn) {
        delete messageSNs[sn];
        console.log('rush-reply', args);
        ipcRenderer.removeListener('rush-reply', reply);
        resolve(args);
      }
    };
    ipcRenderer.on('rush-reply', reply);
  }).then(resp => resp);
}



