/* eslint-disable camelcase */
// @flow
import { isNil } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { call, put, all } from 'redux-saga/effects';
import Device from '../Device';
import type {
  tIOContact,
  tIOData,
  tIODirection,
  tIOPort,
  tIOChange
} from './type';
import { CommonLog } from '../../../../common/utils';
import { ioDirection, ioTriggerMode } from './constants';
import { ioSetApi, ioContactApi, ioStatusApi } from '../../../../api/io';
import type { IIOModule } from './interface/IIOModule';

export default class ClsIOModule extends Device implements IIOModule {
  _data: tIOData = { input: '', output: '' };

  _ports: Array<tIOPort> = [];

  _maxInputs: number = 0;

  _maxOutputs: number = 0;

  constructor(
    name: string,
    serialNumber: string,
    config: { input_num: number, output_num: number },
    // eslint-disable-next-line no-unused-vars,flowtype/no-weak-types
    ...rest: Array<any>
  ) {
    super(name, serialNumber);
    const { input_num, output_num } = config;
    this._maxInputs = input_num;
    this._maxOutputs = output_num;
    for (let i = 0; i < this._maxInputs; i += 1) {
      this._ports.push({
        direction: ioDirection.input,
        idx: i
      });
    }
    for (let i = 0; i < this._maxOutputs; i += 1) {
      this._ports.push({
        direction: ioDirection.output,
        idx: i
      });
    }
    /* eslint-disable flowtype/no-weak-types */
    (this: any).setIO = this.setIO.bind(this);
    (this: any).openIO = this.openIO.bind(this);
    (this: any).closeIO = this.closeIO.bind(this);
    (this: any).getStatus = this.getStatus.bind(this);
    (this: any).ioContact = this.ioContact.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  getPort(direction: tIODirection, idx: number): ?tIOPort {
    return this._ports.find(p => p.direction === direction && p.idx === idx);
  }

  hasPort(port: tIOPort): ?boolean {
    return this._ports.findIndex(p => p === port) >= 0;
  }

  _checkValidateIdx(idx: number, ioType: tIODirection): boolean {
    switch (ioType) {
      case 'input':
        return idx <= this._maxInputs;
      case 'output':
        return idx <= this._maxOutputs;
      default:
        return false;
    }
  }

  static bitString2Boolean(bit: string): boolean {
    switch (bit) {
      case '1':
        return true;
      case '0':
        return false;
      default:
        return false;
    }
  }

  _storeDataField(newData: tIOContact): void {
    const { direction, contact } = newData;
    if (typeof contact !== 'string') {
      CommonLog.lError(`IO Data Must Be String!!!!`);
      return;
    }
    this._data[direction] = contact;
  }

  _getIOChanges(newData: tIOContact): Array<tIOChange> {
    const changes: Array<tIOChange> = [];
    const { direction, contact } = newData;
    const eles: string = this._data?.[direction];
    if (!eles) {
      return changes;
    }
    const startIdx = direction === ioDirection.input ? 0 : this._maxInputs;

    [...contact].forEach((newVal, idx) => {
      if (isNil(newVal)) {
        return;
      }
      if (!this._checkValidateIdx(idx, direction)) {
        return;
      }
      if (!eles[idx]) {
        return;
      }
      const newBool = ClsIOModule.bitString2Boolean(newVal);
      const oldBool = ClsIOModule.bitString2Boolean(eles[idx]);

      if (newBool !== oldBool) {
        changes.push({
          port: this._ports[startIdx + idx],
          triggerMode: ioTriggerMode.change
        });
        if (newBool) {
          changes.push({
            port: this._ports[startIdx + idx],
            triggerMode: ioTriggerMode.rising
          });
        } else {
          changes.push({
            port: this._ports[startIdx + idx],
            triggerMode: ioTriggerMode.falling
          });
        }
      }
    });

    return changes;
  }

  *doDispatch(data: tIOContact): Saga<void> {
    try {
      const ret = this.doValidate(data.contact);
      if (!ret) {
        return;
      }

      this._storeDataField(data);
      const inputs = this._getIOChanges(data).map(c => ({
        data: {
          port: c.port,
          triggerMode: c.triggerMode
        },
        source: this.Name,
        time: new Date()
      }));

      let actions = [];
      inputs.forEach(i => {
        actions = [...actions, ...this._inputListener.check(i)];
      });

      yield all(actions.map(a => put(a)));
    } catch (e) {
      CommonLog.lError(e, {
        at: 'doDispatch',
        data
      });
    }
  }

  *setIO(port: tIOPort, value: boolean): Saga<void> {
    try {
      const status = value ? 1 : 0;
      // eslint-disable-next-line flowtype/no-weak-types
      yield call((ioSetApi: Function), this.serialNumber, port.idx, status);
    } catch (e) {
      CommonLog.lError(e, { at: 'setIO' });
    }
  }

  *openIO(port: tIOPort | Array<tIOPort>): Saga<void> {
    try {
      // eslint-disable-next-line flowtype/no-weak-types
      const openPort = p =>
        call((ioSetApi: Function), this.serialNumber, p.idx, 1);
      if (port instanceof Array) {
        yield all(port.map(p => openPort(p)));
      } else {
        yield openPort(port);
      }
    } catch (e) {
      CommonLog.lError(e, { at: 'openIO' });
    }
  }

  *closeIO(port: tIOPort | Array<tIOPort>): Saga<void> {
    try {
      // eslint-disable-next-line flowtype/no-weak-types
      const closePort = p =>
        call((ioSetApi: Function), this.serialNumber, p.idx, 0);
      if (port instanceof Array) {
        yield all(port.map(p => closePort(p)));
      } else {
        yield closePort(port);
      }
    } catch (e) {
      CommonLog.lError(e, { at: 'closeIO' });
    }
  }

  *getStatus(): Saga<void> {
    try {
      // eslint-disable-next-line flowtype/no-weak-types
      yield call((ioStatusApi: Function), this.serialNumber);
    } catch (e) {
      CommonLog.lError(e, { at: 'getStatus' });
    }
  }

  *ioContact(): Saga<void> {
    try {
      // eslint-disable-next-line flowtype/no-weak-types
      yield call((ioContactApi: Function), this.serialNumber);
    } catch (e) {
      CommonLog.lError(e, { at: 'ioContact' });
    }
  }
}
