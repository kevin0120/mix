import OWebSocket from 'ws';

const WebSocket = require('@oznu/ws-connect');
const { ipcMain } = require('electron');

let ws = null;

let mainWindow = null;

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
  if (messageSNs[data.sn]) { // is reply
    messageSNs[data.sn](data);
    delete messageSNs[data.sn];
  }
  if (mainWindow) {
    mainWindow.send('rush-message', resp);
  }
};

const rushReply = event => resp => {
  console.log('ws reply data: ', resp);
  event.reply('rush-reply', resp, resp.sn);
};

function startListenSend() {
  const replyWSError = onWSMessage();
  ipcMain.on('rush-send', (event, { data, timeout, sn }) => {
    messageSNs[sn] = rushReply(event);
    if (ws && !ws.closed && ws.ws.readyState === OWebSocket.OPEN) {
      const msg = {
        sn,
        ...data
      };
      console.log('sending:');
      console.log(msg);
      ws.sendJson(msg, err => {
        if (err && ws) {
          console.error(err);
          if (!messageSNs[sn]) {
            return;
          }
          replyWSError({
            data: {
              result: -2,
              msg: `error when sending message`
            },
            sn
          });
        }
      });

      setTimeout(() => {
        if (!messageSNs[sn]) {
          return;
        }
        // eslint-disable-next-line no-param-reassign
        replyWSError({
          data: {
            result: -1,
            msg: `rush send timeout`
          },
          sn
        });
      }, timeout);
    } else {
      replyWSError({
        data: {
          result: -404,
          msg: `cannot send message to rush now, rush is not connected`
        },
        sn
      });
    }
  });
}

export function init(url, hmiSN, window) {
  mainWindow = window;
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
      ws.on('websocket-status', (...args) => {
        window.send('rush-status', ...args);
      });
    } else if (!ws.closed && ws.ws.readyState === OWebSocket.OPEN) {
      window.send('rush-open');
    }
  });
  startListenSend();
}
