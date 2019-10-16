// @flow

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

export type tDeviceSN = string;
