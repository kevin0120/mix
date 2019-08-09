// @flow

import { isNil } from 'lodash-es';

type tArrayDevices= {
  // eslint-disable-next-line flowtype/no-weak-types
  [type: string]: Array<any>
};

export const symController = 'Controller';
export const symScanner = 'Scanner';
export const symTool = 'ToolScrew';
export const symReader = 'Reader';
export const symIO = 'IO';


/* eslint-disable flowtype/no-weak-types */
export function AppendNewDevices(symbol: string, clsObj: any) {
  const cl = gDevices?.[symbol];
  if (!(cl instanceof Array)) {
    gDevices[symbol] = [clsObj];
  } else {
    cl.push(clsObj);
  }
  console.log(gDevices);
}

const gDevices: tArrayDevices = {};
// const gExternalSystems: tArrayDevices<any> = {};
/* eslint-enable flowtype/no-weak-types */
