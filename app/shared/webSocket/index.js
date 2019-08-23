// import { CommonLog } from '../../common/utils';

const WebSocket = require('@oznu/ws-connect');


let ws = null;
let open = false;

export function getWSClient() {
  return ws;
}

export function setWSClient(newWS) {
  ws = newWS;
}

export function init(url,hmiSN) {
  ws = new WebSocket(url, {
    reconnectInterval: 3000,
    options:
      {
        maxPayload: 200 * 1024 * 1024
      }
  });
  if(ws){
    ws.on('open', () => {
      ws.sendJson({
        type: 'WS_REG',
        sn: 0,
        data: {
          hmi_sn: hmiSN
        }
      }, err => {
        if (err && ws) {
          console.error(err);
          // CommonLog.lError(err);
        }
      })
    });
  }
}

export function didOpen() {
  if(open){
    open=false;
    return true;
  }
  return false;
}
