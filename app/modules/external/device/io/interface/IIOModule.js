// @flow

import type { Saga } from 'redux-saga';
import type { IDevice } from '../../IDevice';
import type { tIOChange, tIOContact, tIODirection, tIOListener, tIOPort, tIOTriggerMode } from '../type';
import type { tAction } from '../../../../typeDef';

export interface IIOModule extends IDevice {
  getPort(direction: tIODirection, idx: number): ?tIOPort,

  closeIO(port: tIOPort | Array<tIOPort>): Saga<void>,

  hasPort(port: tIOPort): ?boolean,

  _checkValidateIdx(idx: number, ioType: tIODirection): boolean,

  _storeDataField(newData: tIOContact): void,

  _getIOChanges(newData: tIOContact): Array<tIOChange>,

  _doHandleIOData(newData: tIOContact, ...actionParams: any): Array<tAction<any, any>>,

  setIO(port: tIOPort, value: boolean): Saga<void>,

  openIO(port: tIOPort | Array<tIOPort>): Saga<void>,

  getStatus(): Saga<void>,

  ioContact(): Saga<void>,

  addListener(
    port: tIOPort,
    triggerMode: tIOTriggerMode,
    // eslint-disable-next-line flowtype/no-weak-types
    dispatcher: (...args: any) => tAction<any, any>
  ): tIOListener,

  removeListener(listener: tIOListener): Array<tIOListener>
}