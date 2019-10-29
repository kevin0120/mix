// @flow
import type { Saga } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { CommonLog } from '../../../common/utils';
import type Device from './Device';
import ClsController from './controller/model';
import ClsScanner from './scanner/ClsScanner';
import ClsScrewTool from './tools/ClsScrewTool';
import ClsReader from './reader/ClsReader';
import ClsIOModule from './io/ClsIOModule';
import type { tDeviceSN, tArrayDevices, tDeviceType } from './typeDef';
import { deviceStatusApi } from '../../../api/device';
import healthzActions from '../../healthz/action';
import type { tRushData } from '../../rush/type';
import type { IDevice } from './IDevice';
import type { ICommonExternalEntity } from '../ICommonExternalEntity';
import { deviceType, status2Healthz } from './constants';

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
  // eslint-disable-next-line flowtype/no-weak-types
  config: { [string]: any },
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
        device.appendChildren((child: ICommonExternalEntity));
      } else {
        lostChildren[cSN] = device;
      }
    });

    return device;
  } catch (e) {
    CommonLog.lError(e);
  }
}

export function getDevice(sn: tDeviceSN): IDevice {
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

// eslint-disable-next-line flowtype/no-weak-types
export function* deviceStatus(data: tRushData<any, any>): Saga<void> {
  try {
    if (!(data?.data instanceof Array)) {
      return;
    }
    data.data.forEach(d => {
      const { sn, type, children, status, dData, config } = d;
      let dv = getDevice(sn);
      // try make a new device if dv doesn't exist
      if (!dv) {
        dv = newDevice(type, `${type}-${sn}`, sn, config, dData, children);
      }

      // if dv exists, set its Healthz status
      if (dv) {
        dv.Healthz = status2Healthz[status] || false;
        return;
      }
      CommonLog.lError(`invalid device: ${sn}`, { at: 'deviceStatus' });
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
