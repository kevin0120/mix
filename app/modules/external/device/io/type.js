// @flow
import type { AnyAction } from '../type';

export const sOff = 0;
export const sOn = 1;
export const sBlinkOff = 10;
export const sBlinkOn = 11;

export type tIOWSMsgType = 'WS_IO_STATUS' | 'WS_IO_CONTACT' | 'WS_IO_SET';

export const ioDirection = {
  input: 'input',
  output: 'output'
};

export type tIODirection = $Keys<typeof ioDirection>;

export type tIOContact = {
  +sn: string,
  +direction: tIODirection,
  +contact: string // 位串
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

export const DefaultInput = '00000000';
export const DefaultOutput = '00000000';

// 上升沿，下降沿，双向(toggle)

export const ioTriggerMode = {
  rising: 'rising',
  falling: 'falling',
  // high:'high',
  // low:'low',
  change: 'change'
};
// export type tIOTriggerMode = 'Rising' | 'Falling' | 'Bidirectional';
export type tIOTriggerMode = $Keys<typeof ioTriggerMode>;

export type iIODataField = {
  data: boolean,
  triggerMode: tIOTriggerMode,
  // eslint-disable-next-line flowtype/no-weak-types
  action: (mode: tIOTriggerMode, ...args: any) => AnyAction
};

// IO数据字段，key代表的是哪一位， value代表开或者关和相关的action
export interface iIODataFieldObj {
  [key: number]: iIODataField;
}

export type tIOData = {
  input: string,
  output: string
};

export type tIOListener = {
  port: tIOPort,
  triggerMode: tIOTriggerMode,

  // eslint-disable-next-line flowtype/no-weak-types
  dispatcher: (...args: any) => AnyAction
};

export type tIOPort = {
  direction: tIODirection,
  idx: number
};

export type tIOChange = {
  port: tIOPort,
  triggerMode: tIOTriggerMode
};

export default IO_FUNCTION;
