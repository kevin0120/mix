/* eslint-disable no-bitwise,no-case-declarations */
// @flow
import { put, call, fork } from 'redux-saga/effects';
import { channel } from 'redux-saga';
import { loginRequest, logoutRequest } from '../user/action';
import { watch } from '../indexSaga';
import {CARD_AUTH} from './action';


// const pcscLite = require('pcsclite');
const pcscLite = null;

let pcscd = null;
let reader = null;

let pcscdChannel = null;
let readerChannel = null;

const readingConfig={
  buffer:Buffer.from([0xff, 0xb2, 0x00, 0x00, 0x06]),
  length:40     // Max expected length of the response
};

const CARD_STATE = {
  REMOVED: 'CARD_STATE_REMOVED',
  INSERTED: 'CARD_STATE_INSERTED'
};

const PCSCD_CHANNEL = {
  READER: 'PCSCD_CHANNEL_READER',
  ERROR: 'PCSCD_CHANNEL_ERROR'
};

const READER_CHANNEL = {
  ERROR: 'READER_CHANNEL_ERROR',
  STATUS: 'READER_CHANNEL_STATUS',
  END: 'READER_CHANNEL_END',
  CONNECT: 'READER_CHANNEL_CONNECT',
  RECEIVE: 'READER_CHANNEL_RECEIVE'
};

const readerWorkers = {
  [READER_CHANNEL.ERROR]: onCardError,
  [READER_CHANNEL.STATUS]: readerStatusChange,
  [READER_CHANNEL.CONNECT]: readerCardConnect,
  [READER_CHANNEL.END]: readerEnd,
  [READER_CHANNEL.RECEIVE]: readerReceive
};
const pcscdWorkers = {
  [PCSCD_CHANNEL.ERROR]: onCardError,
  [PCSCD_CHANNEL.READER]: pcscdReaderConnect
};

export function* cardAuthFlow() {
  try {
    pcscd = yield call(pcscLite);
    yield call(pcscdListener);
  } catch (e) {
    console.error(e);
  } finally {
    console.log('cardAuthFlow finished');
  }
}

export function* onCardRemoved() {
  try {
    yield put(logoutRequest('112233')); // 默认uuid 112233 登出
  } catch (e) {
    console.error(e);
  }
}

function* pcscdListener() {
  try {
    pcscdChannel = yield call(creatPcscdChannel);
    yield call(watch(pcscdWorkers,pcscdChannel));
  } catch (err) {
    console.error(err);
  }
}

function* pcscdReaderConnect(action) {
  try {
    // eslint-disable-next-line prefer-destructuring
    reader = action.reader;
    readerChannel = yield call(creatReaderChannel);
    yield fork(watch(readerWorkers,readerChannel));
  } catch (e) {
    console.error(e);
  }
}

function onCardError(action) {
  try {
    console.error(action.error);
  } catch (e) {
    console.error(e);
  }
}

const parseBuffer = data => {
  const uuid = '112233';
  console.log(data);
  console.log(data.toString('utf8', 0, -2));
  // uuid=data.toString('utf8',0,-2);
  // TODO

  return uuid;
};

function* readerReceive(action) {
  try {
    const uuid = parseBuffer(action.data);
    yield put(doUserAuth(uuid));
  } catch (e) {
    console.error(e);
  }
}

const cardState = (changes, status) => {
  if (changes & reader.SCARD_STATE_EMPTY &&
    status.state & reader.SCARD_STATE_EMPTY) {
    return CARD_STATE.REMOVED;
  }
  if (changes & reader.SCARD_STATE_PRESENT &&
    status.state & reader.SCARD_STATE_PRESENT) {
    return CARD_STATE.INSERTED;
  }
};

function* readerStatusChange(action) {
  try {
    const { status } = action;
    /* check what has changed */
    const changes = reader.state ^ status.state;
    if (changes) {
      yield call(cardAuthReaderStateChange, cardState(changes, status));
    }
  } catch (e) {
    console.error(e);
  }
}

function readerCardConnect(action) {
  try {
    // read data
    reader.transmit(
      readingConfig.buffer,
      readingConfig.length,
      action.protocol,
      (error, data) => {
        if (error) {
          readerChannel.put({
            type: READER_CHANNEL.ERROR,
            error
          });
        } else {
          readerChannel.put({
            type: READER_CHANNEL.RECEIVE,
            data
          });
        }
      }
    );
  } catch (e) {
    console.error(e);
  }
}

function* readerEnd() {
  try {
    console.log('Reader removed');
    yield call(onCardRemoved);
  } catch (e) {
    console.error(e);
  }
}

// 授权卡状态变化
function* cardAuthReaderStateChange(state) {
  try {
    if (state === CARD_STATE.REMOVED) {
      yield call(onCardRemoved);
      reader.disconnect(reader.SCARD_LEAVE_CARD, error => {
        if (error) {
          readerChannel.put({
            type: READER_CHANNEL.ERROR,
            error
          });
        }
      });
    } else {
      yield put({ type: CARD_AUTH.CARD.INSERTED });
      reader.connect(
        { share_mode: reader.SCARD_SHARE_SHARED },
        (error, protocol) => {
          if (error) {
            readerChannel.put({
              type: READER_CHANNEL.ERROR,
              error
            });
          } else {
            readerChannel.put({
              type: READER_CHANNEL.CONNECT,
              protocol
            });
          }
        }
      );
    }
  } catch (e) {
    console.error(e);
  }
}

function creatPcscdChannel() {
  const c = channel();
  pcscd.on('reader', r => {
    c.put({
      type: PCSCD_CHANNEL.READER,
      reader: r
    });
  });

  pcscd.on('error', err => {
    c.put({
      type: PCSCD_CHANNEL.ERROR,
      error: err
    });
  });
  return c;
}

function creatReaderChannel() {
  const c = channel();
  reader.on('error', err => {
    c.put({
      type: READER_CHANNEL.ERROR,
      error: err
    });
  });

  reader.on('status', status => {
    c.put({
      type: READER_CHANNEL.STATUS,
      status
    });
  });

  reader.on('end', () => {
    c.put({
      type: READER_CHANNEL.END
    });
  });
  return c;
}
