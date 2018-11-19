import OWebSocket from 'ws';
import { RUSH } from './actionTypes';

let ws = null;

const WebSocket = require('@oznu/ws-connect');

export function NewResults(results) {
  return {
    type: RUSH.NEW_RESULTS,
    data: results
  };
}

export function stopWebsocket() {
  if (
    ws.ws.readyState === OWebSocket.OPEN ||
    ws.ws.readyState === OWebSocket.CONNECTING
  ) {
    ws.close();
  }
  ws = null;
}

export function initRush(dispatch, httpConn, hmiSN) {
  const conn = httpConn.slice(7);

  if (conn.length < 7) {
    return;
  }
  const wsURL = `ws://${conn}/rush/v1/ws`;

  if (ws) {
    stopWebsocket();
  }

  ws = new WebSocket(wsURL);
  ws.on('open', () => {
    // reg msg
    ws.sendJson({ hmi_sn: hmiSN }, err => {
      if (err) {
        ws.close();
      }
    });
  });

  ws.on('close', (code, reason) => {
    console.log(
      `websocket disconnected. retry in 1s code: ${code}, reason: ${reason}`
    );
  });

  ws.on('error', () => {
    console.log('websocket error. reconnect after 1s');
  });
  ws.on('ping', () => {
    console.log(' receive ping Msg');
  });
  ws.on('pong', () => {
    console.log(' receive pong Msg');
  });

  ws.on('message', dataRaw => {
    const dataArray = dataRaw.split(';');

    const event = dataArray[0].split(':').slice(-1)[0];

    const data = dataArray.slice(-1);
    const json = JSON.parse(data);

    // const currentState = getState();

    switch (event) {
      case 'job':
        break;
      case 'io':
        break;
      case 'result':
        dispatch(NewResults(json));
        break;
      case 'scanner':
        break;
      default:
        console.log('no define event');
    }
  });
}
