import { eventChannel, channel } from 'redux-saga';
import { take, put, call, fork, all, select, cancel, cancelled } from 'redux-saga/effects';
import modbus from 'jsmodbus';
import { cloneDeep } from 'lodash';
import { IO, OPERATION } from '../actions/actionTypes';
import { IO_FUNCTION } from '../reducers/io';
import { continueOperation } from './operation';

import { setHealthzCheck } from '../actions/healthCheck';
import userConfigs from '../shared/config';
import { setNewNotification } from '../actions/notification';

import { Interval } from './timer';

const Reconnect = require('node-net-reconnect');
const net = require('net');

const lodash = require('lodash');

const IO_CHANNEL = {
  KEY_MONITOR: {
    ON_TICK: 'IO_CHANNEL_KEY_MONITOR_ON_TICK',
    CLOSE: 'IO_CHANNEL_KEY_MONITOR_CLOSE'
  },
  CLOSE_ALL: 'IO_CHANNEL_CLOSE_ALL',
  SENDER: {
    ON_TICK: 'IO_CHANNEL_SENDER_ON_TICK',
    CLOSE: 'IO_CHANNEL_SENDER_CLOSE'
  },
  RECON: {
    CLOSE: 'IO_CHANNEL_RECON_CLOSE'
  },
  FUNCTION: 'IO_CHANNEL_FUNCTION'
};

const CLIENT_CHANNEL = {
  CONNECT: 'CLIENT_CHANNEL_CONNECT',
  END: 'CLIENT_CHANNEL_END',
  CLOSE: 'CLIENT_CHANNEL_CLOSE',
  ERROR: 'CLIENT_CHANNEL_ERROR'
};
// let fail = false;


let oWhite = 0;
let oYellow = 1;
let oGreen = 2;
let oRed = 3;
let oBeep = 4;

const io = {
  channel: null,
  client: null,
  recon: null,
  currentKeyStatus: null,
  modbusClient: null,
  keyMonitorTask: null,
  senderTask: null,
  runningTask: null
};
let iResetKey = 0;
let iBypass = 1;
let iModeSelect = 2;

let iPrevWhite = oWhite;
let iPrevYellow = oYellow;
let iPrevGreen = oGreen;
let iPrevRed = oRed;
let iPrevBeep = oBeep;

export const sOff = 0;
export const sOn = 1;
export const sBlinkOff = 10;
export const sBlinkOn = 11;

let ioStatus = [sOn, sOn, sOn, sOn];

let timeStamp = null;
const byPassTimeout = 3;


export function* watchIO() {
  io.channel = yield call(channel);
  while (true) {
    const action = yield take([IO.INIT, IO.TEST, IO.RESET]);
    switch (action.type) {
      // case IO.FUNCTION:
      //   yield call(handleIOFunction, action.data);
      //   break;
      case IO.INIT:
        yield fork(initIOModbus);
        break;
      // case IO.TEST:
      //   yield fork(testIO, action.io, action.idx);
      //   break;
      case IO.RESET:
        yield fork(resetIO, action.modbusConfig);
        break;
      default:
        break;
    }
  }
}

function* watchIOChannel() {
  try {
    while (true) {
      const action = yield take(io.channel);
      switch (action.type) {
        case IO_CHANNEL.FUNCTION:
          yield call(handleIOFunction, action.data);
          break;
        case IO_CHANNEL.CLOSE_ALL:
          if (io.keyMonitorTask) {
            yield cancel(io.keyMonitorTask);
            io.keyMonitorTask = null;
          }
          if (io.senderTask) {
            yield cancel(io.senderTask);
            io.senderTask = null;
          }
          if (io.recon) {
            io.recon.end();
            io.recon = null;
          }
          if (io.client) {
            io.client.destroy();
          }
          if (io.runningTask) {
            yield cancel(io.runningTask);
            io.runningTask = null;
          }
          break;
        case IO_CHANNEL.SENDER.CLOSE:
          if (io.senderTask) {
            yield cancel(io.senderTask);
            io.senderTask = null;
          }
          break;
        case IO_CHANNEL.KEY_MONITOR.CLOSE:
          if (io.keyMonitorTask) {
            yield cancel(io.keyMonitorTask);
            io.keyMonitorTask = null;
          }
          break;
        default:
          break;
      }
    }
  } finally {
    console.log('watchIOChannel finished');
  }

}

function* handleIOFunction(data) {
  switch (data) {
    case IO_FUNCTION.IN.RESET: {
      // 复位

      yield call(continueOperation);
      break;
    }
    case IO_FUNCTION.IN.BYPASS: {
      // 强制放行

      yield put({ type: OPERATION.FINISHED, data: [] });
      break;
    }
    case IO_FUNCTION.IN.MODE_SELECT: {
      // 模式选择

      break;
    }

    default:
      break;
  }
}

function* runIO() {
  try {
    io.channel.flush(()=>{});
    yield all([
      call(ioClientListener),
      call(watchIOChannel)
    ]);
  } finally {
    console.log('io running finished');
  }
}

function* initIOModbus() {
  const state = yield select();


  yield put(io.channel, { type: IO_CHANNEL.CLOSE_ALL });

  const modbusConfig = state.setting.page.modbus;
  initModBusIO(modbusConfig);

  const modbusIOUrl = state.connections.io;
  const kvs = modbusIOUrl.split('://');
  const hostPorts = kvs[1].split(':');
  const host = hostPorts[0];
  const port = parseInt(hostPorts[1].split('/')[0], 10);

  const options = {
    host,
    port,
    retryTime: 1000, // 1s for every retry
    retryAlways: true // retry even if the connection was closed on purpose
  };


  io.client = new net.Socket();
  io.recon = new Reconnect(io.client, options);
  io.client.setTimeout(1000);
  io.modbusClient = new modbus.client.TCP(io.client);



  try {
    io.client.connect(options);
  } catch (error) {
    io.client = null;
    console.log(error);
    return;
  }
  io.runningTask = yield fork(runIO);

}

function ioClientChannel() {
  return eventChannel(emit => {

    io.client.on('connect', () => {
      emit({ type: CLIENT_CHANNEL.CONNECT });
    });

    io.client.on('end', () => {
      emit({ type: CLIENT_CHANNEL.END });
    });

    io.client.on('close', () => {
      emit({ type: CLIENT_CHANNEL.CLOSE });
    });

    io.client.on('error', () => {
      emit({ type: CLIENT_CHANNEL.ERROR });
    });

    return () => {
    };
  });
}

function* ioClientListener() {
  const chan = yield call(ioClientChannel);
  try {
    while (true) {
      const chanAction = yield take(chan);
      switch (chanAction.type) {
        case CLIENT_CHANNEL.CONNECT:
          yield call(ioClientOnConnect);
          break;
        case CLIENT_CHANNEL.END:
          io.client.end();
          yield put(setHealthzCheck('modbus', false));
          break;
        case CLIENT_CHANNEL.CLOSE:
        case CLIENT_CHANNEL.ERROR:
          yield put(setHealthzCheck('modbus', false));
          break;
        default:
          break;
      }
    }
  }
  finally {
    console.log('ioClientListener finished');
  }
}

function* ioClientOnConnect() {
  yield put(setHealthzCheck('modbus', true));

  io.keyMonitorTask = yield fork(Interval, 500, {
    onTick: [keyMonitorOnTick]
  });
  io.senderTask = yield fork(Interval, 500, {
    onTick: [senderOnTick]
  });
}

function* keyMonitorOnTick() {
  try {
    // set health = true
    const { response, error } = yield call(() => io.modbusClient
      .readDiscreteInputs(iResetKey, 1)
      .then(resp => {
        return ({ response: resp });
      })
      .catch(err => {
        return ({ error: err });
      })
    );
    if (response) {
      const newKeyStatus = response.response.body.valuesAsArray[0];
      if (io.currentKeyStatus !== null && io.currentKeyStatus !== newKeyStatus) {
        if (newKeyStatus === 1) {
          // on
          timeStamp = new Date().getTime();
        } else {
          // off
          const diff = (new Date().getTime() - timeStamp) / 1000;
          if (diff >= byPassTimeout) {
            // 钥匙延迟3秒放行
            if (userConfigs.operationSettings.byPass.type === 'sleep') {
              // yield put({
              //   type: IO.FUNCTION,
              //   data: IO_FUNCTION.IN.BYPASS
              // });
              yield put(io.channel, {
                type: IO_CHANNEL.FUNCTION,
                data: IO_FUNCTION.IN.BYPASS
              });
            }
          } else {
            // 复位动作
            // yield put({
            //   type: IO.FUNCTION,
            //   data: IO_FUNCTION.IN.RESET
            // });
            yield put(io.channel, {
              type: IO_CHANNEL.FUNCTION,
              data: IO_FUNCTION.IN.RESET
            });
          }
        }
      }
// eslint-disable-next-line no-param-reassign
      io.currentKeyStatus = newKeyStatus;
    } else if (error) {
      yield put(setNewNotification('error', error ? error.toString() : 'unknown error'));
      yield cancel(io.keyMonitorTask);
      io.client.destroy();
      yield put(io.channel, IO_CHANNEL.KEY_MONITOR.CLOSE);
    }
  }
  finally {
    if (yield cancelled()) {
      console.log('keyMonitorOnTick canceled');
    }

  }
}

function* senderOnTick() {
  const lights = ioStatus.map(v => {
    if (v === sBlinkOff) {
      return sOff;
    }
    if (v === sBlinkOn) {
      return sOn;
    }
    return v;
  });

  const { response, error } = yield call(() =>
    io.modbusClient
      .writeMultipleCoils(0, lights)
      .then()
      .catch((err) => ({ error: err }))
  );

  if (response) {
  } else if (error) {
    console.log(error);
    io.client.destroy();
    yield put(io.channel, { type: IO_CHANNEL.SENDER.CLOSE });
  }

  ioStatus = ioStatus.map(v => {
    if (v === sBlinkOff) {
      return sBlinkOn;
    }
    if (v === sBlinkOn) {
      return sBlinkOff;
    }
    return v;
  });
}


// export function setFail(f) {
//   fail = f;
// }

function setOutBit(obj) {
  switch (obj.function) {
    case IO_FUNCTION.OUT.LED_RED:
      oRed = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_GREEN:
      oGreen = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_YELLOW:
      oYellow = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_WHITE:
      oWhite = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_BEEP:
      oBeep = obj.bit;
      break;
    default:
      break;
  }
}

function setInBit(obj) {
  switch (obj.function) {
    case IO_FUNCTION.IN.RESET:
      iResetKey = obj.bit;
      break;
    case IO_FUNCTION.IN.BYPASS:
      iBypass = obj.bit;
      break;
    case IO_FUNCTION.IN.MODE_SELECT:
      iModeSelect = obj.bit;
      break;
    default:
      break;
  }
}

function initModBusIO(modbusConfig) {
  const modbusOutConfig = modbusConfig.out;
  const modbusInConfig = modbusConfig.in;
  lodash.forEach(modbusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
}

function resetIO(modbusConfig) {
  const modbusOutConfig = modbusConfig.out;
  const modbusInConfig = modbusConfig.in;
  lodash.forEach(modbusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
  const preIOStatus = cloneDeep(ioStatus);
  ioStatus[oWhite] = preIOStatus[iPrevWhite];
  ioStatus[oYellow] = preIOStatus[iPrevYellow];
  ioStatus[oGreen] = preIOStatus[iPrevGreen];
  ioStatus[oRed] = preIOStatus[iPrevRed];

  iPrevWhite = oWhite;
  iPrevYellow = oYellow;
  iPrevGreen = oGreen;
  iPrevRed = oRed;
  iPrevBeep = oBeep;
}

export function setLedStatusReady() {
  ioStatus[oRed] = sOff;
  ioStatus[oWhite] = sOff;
  ioStatus[oYellow] = sBlinkOn;
  ioStatus[oGreen] = sOff;
}

export function setLedStatusDoing() {
  ioStatus[oRed] = sOff;
  ioStatus[oYellow] = sOn;
  ioStatus[oWhite] = sOn;
  ioStatus[oGreen] = sOff;
}

export function setLedInfo(flag) {
  ioStatus[oWhite] = flag;
}

export function setLedWarning(flag) {
  ioStatus[oGreen] = flag;
}

export function setLedError(flag) {
  ioStatus[oRed] = flag;
}

export function testIO(ioConfig, idx) {
  switch (ioConfig) {
    case 'out': {
      io.modbusClient
        .readCoils(idx, 1)
        .then(resp => {
          const got = resp.response.body.valuesAsArray[0];
          io.modbusClient.writeSingleCoil(idx, got === 0);
          return true;
        })
        .catch(() => 'fail');
      break;
    }
    case 'in': {
      return io.modbusClient.readDiscreteInputs(idx, 1);
    }
    default:
      break;
  }
}
