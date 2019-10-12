import type { tDeviceSN } from './typeDef';

export interface IDevice {
  constructor: (name: string, sn: tDeviceSN, config: Object, data: any)=>IDevice
}