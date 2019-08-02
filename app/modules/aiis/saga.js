import {
  call,
  take,
  takeLatest,
  select,
  put,
  fork,
  cancel
} from 'redux-saga/effects';
import OWebSocket from 'ws';
import { eventChannel } from 'redux-saga';
import { AIIS } from './action';
import { setHealthzCheck } from '../healthzCheck/action';
import { andonNewData } from '../andon/action';
import { setNewNotification } from '../notification/action';

let task = null;

let ws = null;

const WebSocket = require('@oznu/ws-connect');

const AIIS_WS_CHANNEL = {
  OPEN: 'AIIS_WS_CHANNEL_OPEN',
  CLOSE: 'AIIS_WS_CHANNEL_CLOSE',
  ERROR: 'AIIS_WS_CHANNEL_ERROR',
  PING: 'AIIS_WS_CHANNEL_PING',
  PONG: 'AIIS_WS_CHANNEL_PONG',
  MESSAGE: 'AIIS_WS_CHANNEL_MESSAGE'
};

export function* watchAiis() {
  yield takeLatest(AIIS.INIT, initAiis);
}

function* initAiis() {
  try {
    const state = yield select();
    const aiisUrl = state.setting.system.connections.aiis;
    // const hmiSN = state.setting.page.odooConnection.hmiSn.value;
    const hmiSN = state.setting.system.connections.workcenterCode;
    if (task) {
      yield cancel(task);
    }

    if (ws) {
      yield call(stopAiisWebsocket);
    }

    const uris = aiisUrl.split('://');
    if (uris.length > 1) {
      const url = `ws://${uris[1]}/aiis/v1/ws`;
      ws = new WebSocket(url, { reconnectInterval: 3000 });

      task = yield fork(aiisWSListener, hmiSN);
      // yield call(wsOnOpen, hmiSN);
    }
  } catch (e) {
    console.error(e);
  }
}

function aiisWSChannel() {
  return eventChannel(emit => {
    ws.on('open', () => {
      emit({ type: AIIS_WS_CHANNEL.OPEN });
    });

    ws.on('close', (code, reason) => {
      emit({
        type: AIIS_WS_CHANNEL.CLOSE,
        code,
        reason
      });
    });

    ws.on('error', () => {
      emit({ type: AIIS_WS_CHANNEL.ERROR });
    });

    ws.on('ping', () => {
      emit({ type: AIIS_WS_CHANNEL.PING });
    });

    ws.on('pong', () => {
      emit({ type: AIIS_WS_CHANNEL.PONG });
    });

    ws.on('message', dataRaw => {
      emit({
        type: AIIS_WS_CHANNEL.MESSAGE,
        dataRaw
      });
    });

    return () => {
    };
  });
}

function* aiisWSListener(hmiSN) {
  try {
    const chan = yield call(aiisWSChannel);
    while (true) {
      const chanAction = yield take(chan);
      switch (chanAction.type) {
        case AIIS_WS_CHANNEL.OPEN:
          yield put(setHealthzCheck('andon', true));
          yield call(wsOnOpen, hmiSN);
          break;
        case AIIS_WS_CHANNEL.CLOSE:
          yield put(setHealthzCheck('andon', false));
          yield put(setNewNotification('Info', `andon连接状态更新: ${false}`));
          break;
        case AIIS_WS_CHANNEL.ERROR:
          yield put(setHealthzCheck('andon', false));
          yield put(
            setNewNotification('Info', `masterPC连接状态更新: ${false}`)
          );
          break;
        case AIIS_WS_CHANNEL.PING:
          // console.log(' receive ping Msg');
          break;
        case AIIS_WS_CHANNEL.PONG:
          // console.log(' receive pong Msg');
          break;
        case AIIS_WS_CHANNEL.MESSAGE:
          yield call(wsOnMessage, chanAction.dataRaw);
          break;
        default:
          break;
      }
    }
  } catch (e) {
    console.error(e);
  } finally {
    console.log('aiisWSListener finished');
  }
}

function wsOnOpen(hmiSN) {
  ws.sendJson({ hmi_sn: hmiSN }, err => {
    if (err) {
      console.error('aiis ws sendJson error:', err);
      ws.close();
    }
  });
}

function* wsOnMessage(dataRaw) {
  try {
    const dataArray = dataRaw.split(';');

    const json = dataArray.slice(-1);
    const data = JSON.parse(json);

    yield put(andonNewData(data));
  } catch (e) {
    console.error(e);
  }
}

function* stopAiisWebsocket() {
  try {
    if (
      ws.ws.readyState === OWebSocket.OPEN ||
      ws.ws.readyState === OWebSocket.CONNECTING
    ) {
      yield put(setHealthzCheck('andon', false));
      yield put(setNewNotification('Info', `andon连接状态更新: ${false}`));
      ws.close();
    }
    ws = null;
  } catch (e) {
    console.error(e);
  }
}
