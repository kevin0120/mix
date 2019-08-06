// @flow

import { isNil } from "lodash-es";

interface tArrayDevices<t>{
  [type: string]: Array<t>
}

export const symController = 'Controller';
export const symScanner = 'Scanner';
export const symTool = 'ToolScrew';
export const symReader = 'Reader';
export const symIO = 'IO';


/* eslint-disable flowtype/no-weak-types */
export function AppendNewDevices(symbol: string, clsObj: any) {
  const cl = gDevices?.[symbol];
  if (isNil(cl)){
    gDevices[symbol] = [clsObj];
  }else {
    cl.push(clsObj);
  }
}

const gDevices: tArrayDevices<any> = {};
// const gExternalSystems: tArrayDevices<any> = {};
/* eslint-enable flowtype/no-weak-types */
