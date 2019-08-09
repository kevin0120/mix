// @flow

import type { tCommonActionType } from '../../common/type';
import type { tIOContact } from './type';

export const IO_ACTION = {
  DATA_ONCHANGE: 'IO_DATA_ON_CHANGE',
  READ_DATA: 'IO_READ_DATA',
  WRITE_DATA: 'IO_WRITE_DATA',
  RESET: 'IO_RESET'
};

export function onchangeIO(data: ?tIOContact): void | tCommonActionType & {data: ?tIOContact} {
  return {
    type: IO_ACTION.DATA_ONCHANGE,
    data,
  };
}