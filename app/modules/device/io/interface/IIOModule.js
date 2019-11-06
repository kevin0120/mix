// @flow

import type { Saga } from 'redux-saga';
import type { IDevice } from '../../IDevice';
import type { tIOChange, tIOContact, tIODirection, tIOPort } from '../type';

export interface IIOModule extends IDevice {
  getPort(direction: tIODirection, idx: number): ?tIOPort,

  closeIO(port: tIOPort | Array<tIOPort>): Saga<void>,

  hasPort(port: tIOPort): ?boolean,

  _checkValidateIdx(idx: number, ioType: tIODirection): boolean,

  _storeDataField(newData: tIOContact): void,

  _getIOChanges(newData: tIOContact): Array<tIOChange>,

  setIO(port: tIOPort, value: boolean): Saga<void>,

  openIO(port: tIOPort | Array<tIOPort>): Saga<void>,

  getStatus(): Saga<void>,

  ioContact(): Saga<void>
}