// @flow

import type { tCommonActionType } from '../../common/type';
import type { tIOContact } from './type';
import {isNil} from 'lodash-es';

export const IOACTION = {
  DATA_ONCHANGE: 'IO_DATA_ON_CHANGE',
  READ_DATA: 'IO_READ_DATA',
  WRITE_DATA: 'IO_WRITE_DATA',
  RESET: 'IO_RESET'
};

export function onchangeIO(data: ?tIOContact): void | tCommonActionType & {data: ?tIOContact} {
  // if (isNil(data)){
  //   return;
  // }
  return {
    type: IOACTION.DATA_ONCHANGE,
    data,
  };
}

export function readIO(): tCommonActionType {
  return {
    type: IOACTION.READ_DATA
  };
}

// export function resetIO(modbusConfig: any):  {
//   return {
//     type: IO.RESET,
//     modbusConfig
//   };
// }
