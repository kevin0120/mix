// import net from 'net';
import modbus from 'jsmodbus';
import { cloneDeep } from 'lodash';
import userConfigs from '../shared/config';

import { HEALTH, IO } from './actionTypes';
import { IO_FUNCTION } from '../reducers/io';
import { setHealthzCheck } from './healthCheck';
import { setNewNotification } from './notification';

const Reconnect = require('node-net-reconnect');
const net = require('net');

const lodash = require('lodash');

// import {
//   RESET_PAGE_STATUS,
// } from 'actions';
//
// import {
//   IO_TEST_RESPONSE,
// } from 'actions/actionTypes';

let keyMonitorTimer = null;
let senderTimer = null;

let modbusClient = null;
let recon = null;
let client = null;
let fail = false;

let currentKeyStatus = null;

let timeStamp = null;

const byPassTimeout = 3;

let oWhite = 0;
let oYellow = 1;
let oGreen = 2;
let oRed = 3;
let oBeep = 4;

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

export function setFail(f) {
  fail = f;
}

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
    case IO_FUNCTION.OUT.BEEP:
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

function initModBusIO(modusConfig) {
  const moddusOutConfig = modusConfig.out;
  const modbusInConfig = modusConfig.in;
  lodash.forEach(moddusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
}

export function resetIO(modusConfig) {
  const moddusOutConfig = modusConfig.out;
  const modbusInConfig = modusConfig.in;
  lodash.forEach(moddusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
  const preioStatus = cloneDeep(ioStatus);
  ioStatus[oWhite] = preioStatus[iPrevWhite];
  ioStatus[oYellow] = preioStatus[iPrevYellow];
  ioStatus[oGreen] = preioStatus[iPrevGreen];
  ioStatus[oRed] = preioStatus[iPrevRed];

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

export function initIOModbus(dispatch, getState) {
  // 第一步先关闭所有连接
  closeAll();

  const state = getState();

  const modbusConfig = state.setting.page.modbus;
  initModBusIO(modbusConfig);

  client = new net.Socket();

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

  recon = new Reconnect(client, options);
  client.setTimeout(1000);

  modbusClient = new modbus.client.TCP(client);

  try {
    client.connect(options);
  } catch (error) {
    client = null;
    return;
  }

  client.on('connect', () => {
    dispatch(setHealthzCheck('modbus', true));

    keyMonitorTimer = setInterval(() => {
      // set health = true
      modbusClient
        .readDiscreteInputs(iResetKey, 1)
        .then(resp => {
          const newKeyStatus = resp.response.body.valuesAsArray[0];
          if (currentKeyStatus !== null && currentKeyStatus !== newKeyStatus) {
            if (newKeyStatus === 1) {
              // on

              timeStamp = new Date().getTime();
            } else {
              // off

              const diff = (new Date().getTime() - timeStamp) / 1000;
              if (diff >= byPassTimeout) {
                // 钥匙延迟3秒放行
                if (userConfigs.operationSettings.byPass.type === 'sleep') {
                  dispatch({
                    type: IO.FUNCTION,
                    data: IO_FUNCTION.IN.BYPASS
                  });
                }
              } else {
                // 复位动作
                dispatch({
                  type: IO.FUNCTION,
                  data: IO_FUNCTION.IN.RESET
                });
              }
            }
          }
          currentKeyStatus = newKeyStatus;
        })
        .catch(e => {
          dispatch(setNewNotification('error', e.toString()));
          // console.log(error);
          client.destroy();
          clearInterval(keyMonitorTimer);
        });
    }, 500);

    senderTimer = setInterval(() => {
      const lights = ioStatus.map(v => {
        if (v === sBlinkOff) {
          return sOff;
        }
        if (v === sBlinkOn) {
          return sOn;
        }
        return v;
      });
      modbusClient
        .writeMultipleCoils(0, lights)
        .then()
        .catch(() => {
          // console.log(error);
          client.destroy();
          clearInterval(senderTimer);
        });

      ioStatus = ioStatus.map(v => {
        if (v === sBlinkOff) {
          return sBlinkOn;
        }
        if (v === sBlinkOn) {
          return sBlinkOff;
        }
        return v;
      });
    }, 500);
  });

  client.on('end', () => {
    const { healthzStatus } = getState().healthCheckResults;
    if (healthzStatus.modbus.isHealth) {
      // 如果不相等 更新
      dispatch(setHealthzCheck('modbus', false));
    }

    client.end();
  });

  client.on('close', () => {
    const { healthzStatus } = getState().healthCheckResults;
    if (healthzStatus.modbus.isHealth) {
      // 如果不相等 更新
      dispatch(setHealthzCheck('modbus', false));
    }
  });

  client.on('error', () => {
    const { healthzStatus } = getState().healthCheckResults;
    if (healthzStatus.modbus.isHealth) {
      // 如果不相等 更新
      dispatch(setHealthzCheck('modbus', false));
    }
  });
}

export function closeAll() {
  if (keyMonitorTimer) {
    clearInterval(keyMonitorTimer);
  }
  if (senderTimer) {
    clearInterval(senderTimer);
  }
  if (recon) {
    recon.end();
    recon = null;
  }
  if (client) {
    client.destroy();
  }
}

export function testIO(io, idx) {
  switch (io) {
    case 'out': {
      modbusClient
        .readCoils(idx, 1)
        .then(resp => {
          const got = resp.response.body.valuesAsArray[0];
          modbusClient.writeSingleCoil(idx, got === 0);
          return true;
        })
        .catch(() => 'fail');
      break;
    }
    case 'in': {
      return modbusClient.readDiscreteInputs(idx, 1);
    }
    default:
      break;
  }
}
