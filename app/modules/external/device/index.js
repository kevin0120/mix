// @flow
import { CommonLog } from '../../../common/utils';
import type { tRushWebSocketData } from '../../rush/type';
import Device from './Device';
import ClsController from './controller/model';
import ClsScanner from './scanner/model';
import ClsScrewTool from './tools/model';
import ClsReader from './reader/model';
import ClsIOModule from './io/model';
import type { tDeviceSN } from './Device';
import type { Saga } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
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

function newDevice(dt: tDeviceType, name: string, sn: string, config: Object, data, childrenSN: Array<tDeviceSN>) {
  try {
    const device = new sym2Device[dt](name, sn, config, data);
    AppendNewDevices(device);

    // if the new device is a lost child
    if (lostChildren[sn]) {
      const d = lostChildren[sn];
      d.appendChildren(device);
      device.patent = device;
    }

    childrenSN.forEach((cSN) => {
      const child = getDevice(cSN);
      if (child) {
        device.appendChildren(child);
        child.parent = device;
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

export function getDevicesByType(dType:tDeviceType): Array<Device> {
  return [...gDevices].filter((d: Device) => d instanceof sym2Device[dType]);
}

export function AppendNewDevices(deviceObj: Device) {
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
    data.data.forEach((d) => {
      const { sn, type, children, status, data, config } = d;
      let dv = getDevice(sn);
      // try make a new device if dv doesn't exist
      if (!dv) {
        dv = newDevice(type, `${type}${([...gDevices].filter(
          gd => gd instanceof sym2Device[type]
        )?.length || 0) + 1}`, sn, config, data, children);
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
