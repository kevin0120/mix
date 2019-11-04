import OWebSocket from 'ws';

const WebSocket = require('@oznu/ws-connect');
const { ipcMain } = require('electron');

let ws = null;

export function getWSClient() {
  return ws;
}

const messageSNs = {};

function parseData(payload) {
  console.log('ws received data: ', payload);
  const d = /(^[^"]*);(.*)/.exec(payload);
  if (d && d[2]) {
    return (() => {
      try {
        return JSON.parse(d[2]);
      } catch (e) {
        return {};
      }
    })();
  }
  return {};
}

const onWSMessage = dataParser => resp => {
  let data = resp;
  if (dataParser) {
    data = dataParser(resp);
  }
  if (messageSNs[data.sn]) {
    messageSNs[data.sn](data);
    delete messageSNs[data.sn];
  }
};

const rushReply = event => data => {
  console.log('ws reply data: ', data);
  event.reply('rush-reply', data, data.sn);
};

function startListenSend() {
  ipcMain.on('rush-send', (event, { data, timeout, sn }) => {
    messageSNs[sn] = rushReply(event);
    if (ws && !ws.closed && ws.ws.readyState === OWebSocket.OPEN) {
      const msg = {
        sn,
        ...data
      };

      ws.sendJson(msg, err => {
        if (err && ws) {
          console.error(err);
          if (!messageSNs[sn]) {
            return;
          }
          onWSMessage()({
            result: -2,
            msg: `error when sending message`,
            sn
          });
        }
      });

      setTimeout(() => {
        if (!messageSNs[sn]) {
          return;
        }
        // eslint-disable-next-line no-param-reassign
        onWSMessage()({
          result: -1,
          msg: `rush send timeout`,
          sn
        });
      }, timeout);
    } else {
      onWSMessage()({
        result: -404,
        msg: `cannot send message to rush now, rush is not connected`,
        sn
      });
    }
  });
}

export function init(url, hmiSN, window) {
  ipcMain.on('rush', () => {
    const wsMessage = onWSMessage(parseData);
    if (!ws) {
      ws = new WebSocket(url, {
        reconnectInterval: 3000,
        options: {
          maxPayload: 200 * 1024 * 1024
        }
      });
      ws.on('open', () => {
        ws.sendJson(
          {
            type: 'WS_REG',
            sn: 0,
            data: {
              hmi_sn: hmiSN
            }
          },
          err => {
            if (err && ws) {
              console.error(err);
            }
          }
        );
        ws.on('message', wsMessage);
        window.send('rush-open');
      });
      ws.on('close', (...args) => {
        ws.removeListener('message', wsMessage);
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
    } else if (!ws.closed && ws.ws.readyState === OWebSocket.OPEN) {
      window.send('rush-open');
    }
  });
  startListenSend();
}
