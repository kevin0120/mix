import OWebSocket from 'ws';

const WebSocket = require('@oznu/ws-connect');
const { ipcMain } = require('electron');

let ws = null;

export function getWSClient() {
  return ws;
}

const messageSNs = {};

function parseData(payload) {
  const d = /(^[^"]*);(.*)/.exec(payload);
  return JSON.parse(d[2]);
}

const rushReply = (sn, event) => (json) => {
  const data = parseData(json);
  if (data.sn === sn && messageSNs[sn]) {
    event.reply('rush-reply', data, sn);
    ws.removeListener('message', messageSNs[sn]);
    delete messageSNs[sn];
  }
};

function startListenSend() {
  ipcMain.on('rush-send', (event, { data, timeout, sn }) => {
    if (ws && !ws.closed && getWSClient().ws.readyState === OWebSocket.OPEN) {
      const msg = {
        sn,
        ...data
      };

      messageSNs[sn] = rushReply(sn, event);
      ws.on('message', messageSNs[sn]);
      ws.sendJson(msg, (err) => {
        if (err && ws) {
          console.error(err);
          if (!messageSNs[sn]) {
            return;
          }
          ws.removeListener('message', messageSNs[sn]);
          delete messageSNs[sn];
        }
      });

      setTimeout(() => {
        if (!messageSNs[sn]) {
          return;
        }
        // eslint-disable-next-line no-param-reassign
        event.reply('rush-reply', {
          result: -1,
          msg: `rush send timeout`
        }, sn);
        ws.removeListener('message', messageSNs[sn]);
        delete messageSNs[sn];
      }, timeout);

    } else {
      // eslint-disable-next-line no-param-reassign
      event.reply('rush-reply', {
        result: -404,
        msg: `cannot send message to rush now, rush is not connected`
      }, sn);
    }
  });
}

export function init(url, hmiSN, window) {
  ipcMain.once('rush', () => {
    ws = new WebSocket(url, {
      reconnectInterval: 3000,
      options:
        {
          maxPayload: 200 * 1024 * 1024
        }
    });
    if (ws) {
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
          }
        });
        window.send('rush-open');
      });

      ws.on('close', (...args) => {
        window.send('rush-close', ...args);
      });
      ws.on('error', (...args) => {
        window.send('rush-error', ...args);
      });
      ws.on('ping', (...args) => {
        window.send('rush-ping', ...args);
      });
      ws.on('pong', (...args) => {
        window.send('rush-pong', ...args);
      });
      ws.on('message', (...args) => {
        window.send('rush-message', ...args);
      });
      ws.on('websocket-status', (...args) => {
        window.send('rush-status', ...args);
      });
    }
  });
  startListenSend();
}
