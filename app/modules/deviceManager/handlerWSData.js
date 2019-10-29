// @flow
import { put, call, all } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { getDevice, newDevice, getAllDevices } from './devices';
import healthzActions from '../healthz/action';
import type { tRushData } from '../rush/type';
import { CommonLog } from '../../common/utils';
import { status2Healthz } from './constants';
import type { IDevice } from '../external/device/IDevice';
import { makeListener } from '../util';
import type { tAction, tListener } from '../typeDef';

const newDeviceListener = makeListener();

export function* bindNewDeviceListener(
  predicate: IDevice => boolean,
  action: IDevice => tAction<any, any>
): Saga<tListener<IDevice>> {
  const actions = [];
  getAllDevices().forEach(d => {
    if (predicate(d)) {
      actions.push(action(d));
    }
  });
  try {
    yield all(actions.map(a => put(a)));
  } catch (e) {
    CommonLog.lError(e, { at: 'bindNewDeviceListener' });
  }
  return newDeviceListener.add(predicate, action);
}

export function removeNewDeviceListener(listener: tListener<IDevice>) {
  return newDeviceListener.remove(listener);
}

// eslint-disable-next-line flowtype/no-weak-types
export function* deviceStatus(data: tRushData<any, any>): Saga<void> {
  try {
    if (!(data?.data instanceof Array)) {
      return;
    }

    let newDeviceActions = [];
    const setHealthzEffects = [];
    data.data.forEach(d => {
      const { sn, type, children, status, dData, config } = d;
      let dv = getDevice(sn);
      // try make a new device if dv doesn't exist
      if (!dv) {
        dv = newDevice(type, `${type}-${sn}`, sn, config, dData, children);
        newDeviceActions = [
          ...newDeviceActions,
          ...newDeviceListener.check(dv)
        ];
      }

      // if dv exists, set its Healthz status
      if (dv) {
        setHealthzEffects.push(
          // eslint-disable-next-line redux-saga/yield-effects
          call(dv.setHealthz, status2Healthz[status] || false)
        );
        return;
      }
      CommonLog.lError(`invalid device: ${sn}`, { at: 'deviceStatus' });
    });
    yield all(newDeviceActions.map(a => put(a)));
    yield all(setHealthzEffects);

    const status = {};
    getAllDevices().forEach((d: IDevice) => {
      status[d.Name] = d.Healthz;
    });
    yield put(healthzActions.data(status));
  } catch (e) {
    CommonLog.lError(e);
  }
}
