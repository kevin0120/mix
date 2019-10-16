// @flow
import type { Saga } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { CommonLog } from '../../../common/utils';
import type { tRushWebSocketData } from '../../rush/type';
import Device from './Device';
import ClsController from './controller/model';
import ClsScanner from './scanner/ClsScanner';
import ClsScrewTool from './tools/ClsScrewTool';
import ClsReader from './reader/ClsReader';
import ClsIOModule from './io/ClsIOModule';
import type { tDeviceSN } from './typeDef';
import { deviceStatusApi } from '../../../api/device';
import healthzActions from '../../healthz/action';

type tArrayDevices = Set<Device>;

export const deviceType = {
  controller: 'controller',
  scanner: 'scanner',
  tool: 'tool',
  reader: 'reader',
  io: 'io'
};

type tDeviceType = $Keys<typeof deviceType>;

const status2Healthz = {
  online: true,
  offline: false
};

const sym2Device = {
  [deviceType.controller]: ClsController,
  [deviceType.scanner]: ClsScanner,
  [deviceType.tool]: ClsScrewTool,
  [deviceType.reader]: ClsReader,
  [deviceType.io]: ClsIOModule
};

const gDevices: tArrayDevices = new Set([]);

const lostChildren = {};

function newDevice(
  dt: tDeviceType,
  name: string,
  sn: string,
  config: Object,
  data,
  childrenSN: Array<tDeviceSN>
) {
  try {
    const device = new sym2Device[dt](name, sn, config, data);
    appendNewDevices(device);

    // if the new device is a lost child
    if (lostChildren[sn]) {
      const d = lostChildren[sn];
      device.appendChildren(d);
      delete lostChildren[sn];
    }

    // search for the children of the device
    childrenSN.forEach(cSN => {
      const child = getDevice(cSN);
      if (child) {
        device.appendChildren(child);
      } else {
        lostChildren[cSN] = device;
      }
    });

    return device;
  } catch (e) {
    CommonLog.lError(e);
  }
}

export function getDevice(sn: tDeviceSN): Device {
  return [...gDevices].filter((d: Device) => d.serialNumber === sn)?.[0];
}

export function getDevicesByType(dType: tDeviceType): Array<Device> {
  return [...gDevices].filter((d: Device) => d instanceof sym2Device[dType]);
}

export function appendNewDevices(deviceObj: Device) {
  gDevices.add(deviceObj);
}

export function* updateDeviceStatus(): Saga<void> {
  try {
    yield call(deviceStatusApi);
  } catch (e) {
    CommonLog.lError(e);
  }
}

export function* deviceStatus(data: tRushWebSocketData): Saga<void> {
  try {
    if (!(data?.data instanceof Array)) {
      return;
    }
    data.data.forEach(d => {
      const { sn, type, children, status, data, config } = d;
      let dv = getDevice(sn);
      // try make a new device if dv doesn't exist
      if (!dv) {
        dv = newDevice(type, `${type}-${sn}`, sn, config, data, children);
      }

      // if dv exists, set its Healthz status
      if (dv) {
        dv.Healthz = status2Healthz[status] || false;
      }
    });
    const status = {};
    gDevices.forEach((d: Device) => {
      status[d.Name] = d.Healthz;
    });
    yield put(healthzActions.data(status));
  } catch (e) {
    CommonLog.lError(e);
  }
}
