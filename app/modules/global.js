// @flow

import { isNil } from "lodash-es";

interface tArrayDevices<t>{
  [type: string]: Array<t>
}

export const symController = 'Controller';
export const symScanner = 'Scanner';
export const symTool = 'ToolScrew';
export const symReader = 'Reader';


/* eslint-disable flowtype/no-weak-types */
export function AppendNewDevices(symbol: string, clsObj: any) {
  const cl = gDevices?.[symController];
  if (isNil(cl)){
    gDevices[symController] = [clsObj];
  }else {
    cl.push(clsObj);
  }
}

let gDevices: tArrayDevices<any>;
// const gExternalSystems: tArrayDevices<any> = {};
/* eslint-enable flowtype/no-weak-types */
