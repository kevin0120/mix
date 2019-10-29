// @flow
import type {
  tDeviceSN,
  tInputData,
  tInputListener,
  tInputPredicate
} from './typeDef';
import type { ICommonExternalEntity } from '../ICommonExternalEntity';
import type { tAction, tCallable, tListenerObj } from '../../typeDef';

export interface IDevice extends ICommonExternalEntity {
  // eslint-disable-next-line flowtype/no-weak-types
  +constructor: (
    name: string,
    sn: tDeviceSN,
    config: Object,
    data: any
  ) => void;
  +doDispatch: tCallable<tInputData, void>;
  _validator: null | ((data: tInputData) => boolean);
  _inputListener: tListenerObj;
  +serialNumber: ?tDeviceSN;

  addListener(
    predicate: tInputPredicate,
    // eslint-disable-next-line flowtype/no-weak-types
    action: (...args: any) => tAction<any, any>
  ): tInputListener;

  removeListener(listener: tInputListener): Array<tInputListener>;
}
