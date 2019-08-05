// @flow

interface tArrayDevices<t>{
  [type: string]: Array<t>
}

/* eslint-disable flowtype/no-weak-types */
export const gDevices: tArrayDevices<any> = {};
export const gExternalSystems: tArrayDevices<any> = {};
/* eslint-enable flowtype/no-weak-types */
