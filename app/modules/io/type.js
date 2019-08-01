// @flow
import type {AnyAction} from '../../common/type'

export const sOff = 0;
export const sOn = 1;
export const sBlinkOff = 10;
export const sBlinkOn = 11;

export type tIOWSMsgType = 'WS_IO_STATUS' | 'WS_IO_CONTACT' | 'WS_IO_SET';

export type tIOContact =  {
  +sn: string,
  +type: 'input' | 'output',
  +contact: string  // 位串
};

const IO_FUNCTION = {
  IN: {
    RESET: 'RESET',
    BYPASS: 'BYPASS',
    MODE_SELECT: 'MODE_SELECT'
  },
  OUT: {
    LED_WHITE: 'LED_WHITE',
    LED_YELLOW: 'LED_YELLOW',
    LED_GREEN: 'LED_GREEN',
    LED_RED: 'LED_RED',
    BEEP: 'BEEP'
  }
};

/* eslint-disable flowtype/no-weak-types */
export type  iIODataField = {
  +data: boolean,
  +action: (...args: any) => AnyAction
};
/* eslint-enable flowtype/no-weak-types */

// IO数据字段，key代表的是哪一位， value代表开或者关和相关的action
export interface iIODataFieldObj {
  [key: number]: iIODataField
}

export type tIOData = {
  +inputs: iIODataFieldObj,
  +outputs: iIODataFieldObj
};

export default IO_FUNCTION;
