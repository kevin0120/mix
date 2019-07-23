// @flow

import type {tCommonActionType} from '../../common/type';

export const IO = {
  FUNCTION: 'IO_FUNCTION',
  INIT: 'IO_INIT',
  TEST: 'IO_TEST',
  RESET: 'IO_RESET'
};

export function initIO(): tCommonActionType{
  return {
    type: IO.INIT
  };
}

export function resetIO(modbusConfig: any) {
  return {
    type: IO.RESET,
    modbusConfig
  };
}
