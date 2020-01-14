// @flow

import type { tArrayDevices, tDeviceType } from './interface/typeDef';
import type { tDeviceSN } from '../device/typeDef';
import type { IDevice } from '../device/IDevice';
import type { ICommonExternalEntity } from '../external/ICommonExternalEntity';
import ClsController from '../device/controller/ClsController';
import ClsScanner from '../device/scanner/ClsScanner';
import ClsScrewTool from '../device/tools/ClsScrewTool';
import ClsReader from '../device/reader/ClsReader';
import ClsIOModule from '../device/io/ClsIOModule';
import { deviceType } from './constants';

const sym2Device = Object.freeze({
  [deviceType.controller]: ClsController,
  [deviceType.scanner]: ClsScanner,
  [deviceType.tool]: ClsScrewTool,
  [deviceType.reader]: ClsReader,
  [deviceType.io]: ClsIOModule
});

const gDevices: tArrayDevices = new Set([]);

const lostChildren = {};

export function newDevice(
  dt: tDeviceType,
  name: string,
  sn: string,
  // eslint-disable-next-line flowtype/no-weak-types
  config: { [string]: any },
  // eslint-disable-next-line flowtype/no-weak-types
  data: { [string]: any },
  childrenSN: Array<tDeviceSN>
) {
  const device = new sym2Device[dt](name, sn, config, data);
  appendNewDevices(device);

  // if the new device is a lost child
  if (lostChildren[sn]) {
    const parent = getDevice(lostChildren[sn].parent);
    if (parent) {
      parent.appendChildren(device);
      device.setParent(parent);
      delete lostChildren[sn];
    }
  }

  // search for the children of the device
  (Object.keys(childrenSN || {}) || []).forEach(cSN => {
    const child = getDevice(cSN);
    if (child) {
      device.appendChildren((child: ICommonExternalEntity));
      child.setParent(device);
    } else {
      lostChildren[cSN] = {
        device,
        parent: sn
      };
    }
  });

  return device;
}

export function getAllDevices(): Array<IDevice> {
  return [...gDevices];
}

export function getDevice(sn: tDeviceSN): IDevice {
  return [...gDevices].filter((d: IDevice) => d.serialNumber === sn)?.[0];
}

export function getDevicesByType(dType: tDeviceType): Array<IDevice> {
  return [...gDevices].filter((d: IDevice) => d instanceof sym2Device[dType]);
}

export function appendNewDevices(deviceObj: IDevice) {
  gDevices.add(deviceObj);
}
