// @flow
import type { tDeviceSN, tInputData } from './typeDef';
import type { ICommonExternalEntity } from '../ICommonExternalEntity';
import type { tCallable } from '../../typeDef';

export interface IDevice extends ICommonExternalEntity {
  +constructor: (name: string, sn: tDeviceSN, config: Object, data: any)=>void,
  +doDispatch: tCallable<tInputData, void>,
  _validator: null | ((data: tInputData) => boolean),
  +serialNumber: ?tDeviceSN
}