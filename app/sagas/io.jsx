// redux-saga
import { eventChannel, channel, delay } from 'redux-saga';
import {
  take,
  put,
  call,
  fork,
  all,
  select,
  cancel,
  cancelled
} from 'redux-saga/effects';

//
import { cloneDeep } from 'lodash';
import modbus from 'jsmodbus';

// actions
import { IO, OPERATION } from '../actions/actionTypes';
import { setHealthzCheck } from '../actions/healthCheck';
import { setNewNotification } from '../actions/notification';

// reducers
import { IO_FUNCTION } from '../reducers/io';

// sagas
import { continueOperation } from './operation';

// config
import userConfigs from '../shared/config';

const Reconnect = require('node-net-reconnect');
const net = require('net');
const lodash = require('lodash');

const CLIENT_CHANNEL = {
  CONNECT: 'CLIENT_CHANNEL_CONNECT',
  END: 'CLIENT_CHANNEL_END',
  CLOSE: 'CLIENT_CHANNEL_CLOSE',
  ERROR: 'CLIENT_CHANNEL_ERROR'
};

const io = {
  channel: null,
  client: null,
  recon: null,
  currentKeyStatus: null,
  modbusClient: null,
  senderReceiver: null,
  runningTask: null,
  health: false,
  i: {
    resetKey: 0,
    byPass: 1,
    modeSelect: 2
  },
  o: {
    white: 0,
    yellow: 1,
    green: 2,
    red: 3,
    beep: 4
  }
};

export const sOff = 0;
export const sOn = 1;
export const sBlinkOff = 10;
export const sBlinkOn = 11;

let ioStatus = [sOn, sOn, sOn, sOn];

let timeStamp = null;
const byPassTimeout = 3;

export function* watchIO() {
  try {
    io.channel = yield call(channel);
    while (true) {
      const action = yield take([IO.INIT, IO.TEST, IO.RESET]);
      switch (action.type) {
        case IO.INIT:
          yield fork(initIOModbus);
          break;
        // case IO.TEST:
        //   yield fork(testIO, action.io, action.idx);
        //   break;
        case IO.RESET:
          console.log('resetIO');
          yield fork(resetIO, action.modbusConfig);
          break;
        default:
          break;
      }
    }
  } catch (e) {
    console.log(e);
  }
}

function* closeAll() {
  try {
    if (io.senderReceiver) {
      yield cancel(io.senderReceiver);
      io.senderReceiver = null;
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
  } catch (e) {
    console.log(e);
  }
}

function* handleIOFunction(data) {
  try {
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
  } catch (e) {
    console.log(e);
  }
}

function* initIOModbus() {
  try {
    const state = yield select();

    yield call(closeAll);

    const modbusConfig = state.setting.page.modbus;
    setModBusIO(modbusConfig);

    const modbusIOUrl = state.connections.io;
    if (lodash.isNil(modbusIOUrl) || !modbusIOUrl.length) {
      return;
    }
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
    io.runningTask = yield fork(ioClientListener);
  } catch (e) {
    console.log(e);
  }
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

    return () => {};
  });
}

function* ioClientListener() {
  try {
    const chan = yield call(ioClientChannel);
    while (true) {
      const chanAction = yield take(chan);
      switch (chanAction.type) {
        case CLIENT_CHANNEL.CONNECT:
          yield fork(setHealth, true);
          io.senderReceiver = yield fork(senderReceiver);
          break;
        case CLIENT_CHANNEL.END:
          io.client.end();
          yield fork(setHealth, false);
          break;
        case CLIENT_CHANNEL.CLOSE:
        case CLIENT_CHANNEL.ERROR:
          yield fork(setHealth, false);
          break;
        default:
          break;
      }
    }
  } catch (err) {
    console.log(err);
  } finally {
    console.log('ioClientListener finished');
  }
}

function* senderReceiver() {
  try {
    while (true) {
      yield all([call(keyMonitorTask), call(senderTask)]);
      yield delay(500);
    }
  } catch (err) {
    console.log(err);
  } finally {
    console.log('senderReceiver finished');
  }
}

function* keyMonitorTask() {
  try {
    const { response, error } = yield call(() =>
      io.modbusClient
        .readDiscreteInputs(io.i.resetKey, 1)
        .then(resp => ({ response: resp }))
        .catch(err => ({ error: err }))
    );
    if (response) {
      const newKeyStatus = response.response.body.valuesAsArray[0];
      if (
        io.currentKeyStatus !== null &&
        io.currentKeyStatus !== newKeyStatus
      ) {
        if (newKeyStatus === 1) {
          // on
          timeStamp = new Date().getTime();
        } else {
          // off
          const diff = (new Date().getTime() - timeStamp) / 1000;
          if (diff >= byPassTimeout) {
            // 钥匙延迟3秒放行
            if (userConfigs.operationSettings.byPass.type === 'sleep') {
              yield call(handleIOFunction, IO_FUNCTION.IN.BYPASS);
            }
          } else {
            // 复位动作
            yield call(handleIOFunction, IO_FUNCTION.IN.RESET);
          }
        }
      }
      io.currentKeyStatus = newKeyStatus;
    } else if (error) {
      yield put(
        setNewNotification('error', error ? error.toString() : 'unknown error')
      );
      yield cancel(io.senderReceiver);
      io.senderReceiver = null;
      io.client.destroy();
    }
  } finally {
    if (yield cancelled()) {
      console.log('keyMonitorOnTick canceled');
    }
  }
}

function* senderTask() {
  try {
    const lights = ioStatus.map(v => {
      if (v === sBlinkOff) {
        return sOff;
      }
      if (v === sBlinkOn) {
        return sOn;
      }
      return v;
    });

    const { error } = yield call(() =>
      io.modbusClient
        .writeMultipleCoils(0, lights)
        .then()
        .catch(err => ({ error: err }))
    );

    if (error) {
      console.log(error);
      yield cancel(io.senderReceiver);
      io.senderReceiver = null;
      io.client.destroy();
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
  } catch (err) {
    console.log(err);
  }
}

function* setHealth(health) {
  try {
    if (health !== io.health) {
      io.health = health;
      yield put(setHealthzCheck('modbus', health));
    }
  } catch (err) {
    console.log(err);
  }
}

function setOutBit(obj) {
  switch (obj.function) {
    case IO_FUNCTION.OUT.LED_RED:
      io.o.red = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_GREEN:
      io.o.green = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_YELLOW:
      io.o.yellow = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_WHITE:
      io.o.white = obj.bit;
      break;
    case IO_FUNCTION.OUT.BEEP:
      io.o.beep = obj.bit;
      break;
    default:
      break;
  }
}

function setInBit(obj) {
  switch (obj.function) {
    case IO_FUNCTION.IN.RESET:
      io.i.resetKey = obj.bit;
      break;
    case IO_FUNCTION.IN.BYPASS:
      io.i.byPass = obj.bit;
      break;
    case IO_FUNCTION.IN.MODE_SELECT:
      io.i.modeSelect = obj.bit;
      break;
    default:
      break;
  }
}

function setModBusIO(modbusConfig) {
  const modbusOutConfig = modbusConfig.out;
  const modbusInConfig = modbusConfig.in;
  lodash.forEach(modbusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
}

function resetIO(modbusConfig) {
  const preIOStatus = cloneDeep(ioStatus);
  const preO = cloneDeep(io.o);
  const { o } = io;

  setModBusIO(modbusConfig);
  for (const key of o) {
    ioStatus[o[key]] = preIOStatus[preO[key]];
  }
}

export function setLedStatusReady() {
  const { o } = io;
  ioStatus[o.red] = sOff;
  ioStatus[o.white] = sOff;
  ioStatus[o.yellow] = sOff;
  ioStatus[o.green] = sBlinkOn;
}

export function setLedStatusDoing() {
  const { o } = io;
  ioStatus[o.red] = sOff;
  ioStatus[o.yellow] = sOff;
  ioStatus[o.white] = sOn;
  ioStatus[o.green] = sOn;
}

export function setLedInfo(flag) {
  ioStatus[io.o.white] = flag;
}

export function setLedWarning(flag) {
  ioStatus[io.o.green] = flag;
}

export function setLedError(flag) {
  ioStatus[io.o.red] = flag;
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

export function getIBypass() {
  return io.i.byPass;
}

export function getIModeSelect() {
  return io.i.modeSelect;
}
