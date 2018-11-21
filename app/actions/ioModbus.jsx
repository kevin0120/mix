// import net from 'net';
import modbus from 'jsmodbus';
import userConfigs from '../shared/config';

import { HEALTH, IO } from './actionTypes';
import { IO_FUNCTION } from '../reducers/io';

const Reconnect = require('node-net-reconnect');
const net = require('net');

const lodash = require('lodash');

// import { cloneDeep } from 'lodash';

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

let oInfo = 0;
let oDoing = 1;
let oWarning = 2;
let oError = 3;
let oBeep = 4;

let iResetKey = 0;
let iBypass = 1;
let iModeSelect = 2;

// let iPrevRed = iRed;
// let iPrevYellow = iYellow;
// let iPrevGreen = iGreen;
// let iPrevBlue = iBlue;
// let iPrevBeep = iBeep;

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
    case IO_FUNCTION.OUT.LED_ERROR:
      oError = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_WARNING:
      oWarning = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_DOING:
      oDoing = obj.bit;
      break;
    case IO_FUNCTION.OUT.LED_INFO:
      oInfo = obj.bit;
      break;
    case IO_FUNCTION.OUT.BEEP:
      oBeep = obj.bit;
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
  }
}

function initModBusIO(modusConfig) {
  const moddusOutConfig = modusConfig.out;
  const modbusInConfig = modusConfig.in;
  lodash.forEach(moddusOutConfig, setOutBit);
  lodash.forEach(modbusInConfig, setInBit);
}

// export function resetIO(modusConfig) {
//   return dispatch => {
//     const moddusOutConfig = modusConfig.out;
//     const modbusInConfig = modusConfig.in;
//     lodash.forEach(moddusOutConfig, setOutBit);
//     lodash.forEach(modbusInConfig, setInBit);
//     const preioStatus = cloneDeep(ioStatus);
//     ioStatus[iRed] = preioStatus[iPrevRed];
//     ioStatus[iBlue] = preioStatus[iPrevBlue];
//     ioStatus[iGreen] = preioStatus[iPrevGreen];
//     ioStatus[iYellow] = preioStatus[iPrevYellow];
//
//     iPrevRed = iRed;
//     iPrevYellow = iYellow;
//     iPrevGreen = iGreen;
//     iPrevBlue = iBlue;
//     iPrevBeep = iBeep;
//   };
// }

export function setLedStatusReady() {
  ioStatus[oError] = sOff;
  ioStatus[oInfo] = sOff;
  ioStatus[oDoing] = sBlinkOn;
  ioStatus[oWarning] = sOff;
}

export function setLedStatusDoing() {
  ioStatus[oError] = sOff;
  ioStatus[oDoing] = sOn;
  ioStatus[oInfo] = sOn;
  ioStatus[oWarning] = sOff;
}

export function setLedInfo(flag) {
  ioStatus[oInfo] = flag;
}

export function setLedWarning(flag) {
  ioStatus[oWarning] = flag;
}

export function setLedError(flag) {
  ioStatus[oError] = flag;
}

export function initIOModbus(dispatch, state) {
  // 第一步先关闭所有连接
  closeAll();

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
    dispatch({
      type: HEALTH.HEALTH,
      category: HEALTH.IO,
      isHealth: true
    });

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
        .catch(error => {
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
        .catch(error => {
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
    dispatch({
      type: HEALTH.HEALTH,
      category: HEALTH.IO,
      isHealth: false
    });

    client.end();
  });

  client.on('close', () => {
    dispatch({
      type: HEALTH.HEALTH,
      category: HEALTH.IO,
      isHealth: false
    });
  });

  client.on('error', () => {
    dispatch({
      type: HEALTH.HEALTH,
      category: HEALTH.IO,
      isHealth: false
    });
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
