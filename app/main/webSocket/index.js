const WebSocket = require('@oznu/ws-connect');
const { ipcMain } = require('electron');

let ws = null;

export function getWSClient() {
  return ws;
}

export function init(url, hmiSN, window) {
  ipcMain.once('rush', (event) => {
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
}
