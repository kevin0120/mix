import type { Action } from 'redux';

// eslint-disable-next-line flowtype/no-weak-types
export type tInputData = any;

export type tInput = {
  data: tInputData,
  source: string,
  time: Date
};

export type AnyAction = Action & {
  // eslint-disable-next-line flowtype/no-weak-types
  [extraProps: string]: any
};


export type tDeviceNewData = {
  +data: tInputData
};

export type tCommonActionType = {
  +type: string
};
