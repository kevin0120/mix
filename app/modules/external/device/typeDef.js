// @flow

import type { tAction } from '../../typeDef';
import type Device from './Device';
import { deviceType } from './constants';

// eslint-disable-next-line flowtype/no-weak-types
export type tInputData = any;

export type tInput = {
  data: tInputData,
  source: string,
  time: Date
};

export type tDeviceNewData = {
  +data: tInputData
};

export type tInputPredicate = tInput => boolean;

export type tDeviceSN = string;

export type tInputListener = {
  predicate: tInputPredicate,
  // eslint-disable-next-line flowtype/no-weak-types
  action: (...args: any) => tAction<any, any>
};

export type tArrayDevices = Set<Device>;

export type tDeviceType = $Values<typeof deviceType>;
