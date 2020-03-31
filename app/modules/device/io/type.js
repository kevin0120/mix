// @flow
import { IO_WS_TYPES, ioDirection, ioTriggerMode } from './constants';
import type { tAction } from '../../typeDef';

export type tIOWSMsgType = $Values<typeof IO_WS_TYPES>;

export type tIODirection = $Values<typeof ioDirection>;

export type tIOContact = {
  +sn: string,
  +direction: tIODirection,
  +contact: string // 位串
};

export type tIOWSDataContact = {
  src: string, // 设备类型
  sn: string, // 序列号
  inputs?: string,
  outputs?: string
};

// 上升沿，下降沿，双向(toggle)
export type tIOTriggerMode = $Values<typeof ioTriggerMode>;

export type iIODataField = {
  data: boolean,
  triggerMode: tIOTriggerMode,
  // eslint-disable-next-line flowtype/no-weak-types
  action: (mode: tIOTriggerMode, ...args: any) => tAction<any, any>
};

// IO数据字段，key代表的是哪一位， value代表开或者关和相关的action
export interface iIODataFieldObj {
  [key: number]: iIODataField;
}

export type tIOData = {
  input: string,
  output: string
};

export type tIOPort = {
  direction: tIODirection,
  idx: number
};

export type tIOChange = {
  port: tIOPort,
  triggerMode: tIOTriggerMode
};

export type tIOConfig = {
  input_num: number,
  output_num: number
};
