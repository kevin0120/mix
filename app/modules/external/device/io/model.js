// @flow
import { isNil, remove } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { call, fork, put } from 'redux-saga/effects';
import Device from '../Device';
import type { AnyAction } from '../type';
import type {
  tIOContact,
  tIOData,
  tIOTriggerMode,
  tIODirection,
  tIOListener,
  tIOPort,
  tIOChange
} from './type';
import { CommonLog } from '../../../../common/utils';
import { AppendNewDevices, symIO } from '../global';
import { ioDirection, ioTriggerMode } from './type';
import { ioSetApi, ioContactApi, ioStatusApi } from '../../../../api/io';

export default class ClsIOModule extends Device {

  #data: tIOData = { input: '', output: '' };

  #ports: Array<tIOPort> = [];

  #maxInputs: number = 0;

  #maxOutputs: number = 0;

  #listeners: Array<tIOListener> = [];

  #serialNumber: ?string = null;

  get serialNumber() {
    return this.#serialNumber;
  }

  constructor(name: string, serialNumber: string, inputs: number, outputs: number) {
    super(name);
    this.#serialNumber = serialNumber;
    this.#maxInputs = inputs;
    this.#maxOutputs = outputs;
    for (let i = 0; i < inputs; i += 1) {
      this.#ports.push({
        direction: ioDirection.input,
        idx: i
      });
    }
    for (let i = 0; i < outputs; i += 1) {
      this.#ports.push({
        direction: ioDirection.output,
        idx: i
      });
    }
    AppendNewDevices(symIO, this);
    /* eslint-disable flowtype/no-weak-types */
    (this: any).setIO = this.setIO.bind(this);
    (this: any).openIO = this.openIO.bind(this);
    (this: any).closeIO = this.closeIO.bind(this);
    (this: any).getStatus = this.getStatus.bind(this);
    (this: any).ioContact = this.ioContact.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  getPort(direction: tIODirection, idx: number): ?tIOPort {
    return this.#ports.find(p => p.direction === direction && p.idx === idx);
  }

  _checkValidateIdx(idx: number, ioType: tIODirection): boolean {
    switch (ioType) {
      case 'input':
        return idx <= this.#maxInputs;
      case 'output':
        return idx <= this.#maxOutputs;
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
    this.#data[direction] = contact;
  }


  _getIOChanges(newData: tIOContact): Array<tIOChange> {
    const changes: Array<tIOChange> = [];
    const { direction, contact } = newData;
    const eles: string = this.#data?.[direction];
    if (!eles) {
      return changes;
    }
    const startIdx = direction === ioDirection.input ? 0 : this.#maxInputs;

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
          port: this.#ports[startIdx + idx],
          triggerMode: ioTriggerMode.change
        });
        if (newBool) {
          changes.push({
            port: this.#ports[startIdx + idx],
            triggerMode: ioTriggerMode.rising
          });
        } else {
          changes.push({
            port: this.#ports[startIdx + idx],
            triggerMode: ioTriggerMode.falling
          });
        }
      }
    });

    return changes;

  }

  // eslint-disable-next-line flowtype/no-weak-types
  _doHandleIOData(newData: tIOContact, ...actionParams: any): Array<AnyAction> {
    const ret = this.doValidate(newData.contact);
    if (!ret) {
      return [];
    }
    const changes = this._getIOChanges(newData, ...actionParams);
    const matchedListeners = this.#listeners.filter((l) =>
      changes.findIndex(c =>
        c.port === l.port && c.triggerMode === l.triggerMode
      ) >= 0
    );

    this._storeDataField(newData);
    return matchedListeners.map(l => l.dispatcher({
      data: newData,
      source: this.Name,
      time: new Date()
    }));
  }

  // eslint-disable-next-line flowtype/no-weak-types
  set dispatcher(dispatcher: null | (...args: any) => AnyAction) {
    super.dispatcher = null; // 永远设置的是null
  }

  // eslint-disable-next-line class-methods-use-this
  get dispatcher() {
    return null;
  }

  *doDispatch(newData: tIOContact): Saga<void>{
    // CommonLog.Info(`IO Module Please Use doHandleIOData Method`);
    try {
      const actions=this._doHandleIOData(newData)
      if (actions instanceof Array) {
        // eslint-disable-next-line
        for (const a of actions) {
          yield put(a);
        }
      }

    } catch (e) {
      CommonLog.lError(e,{
        at:'doDispatch',
        data:newData
      });
    }
  }

  * setIO(port: tIOPort, value: boolean): Saga<void> {
    try {
      const status = value ? 1 : 0;
      yield call(ioSetApi, this.serialNumber, port.idx, status);
    } catch (e) {
      CommonLog.lError(e, { at: 'setIO' });
    }
  }

  * openIO(port: tIOPort | Array<tIOPort>): Saga<void> {
    try {
      if (port instanceof Array) {
        // eslint-disable-next-line no-restricted-syntax
        for (const p of port) {
          yield fork(ioSetApi, this.serialNumber, p.idx, 1);
        }
      } else {
        yield call(ioSetApi, this.serialNumber, port.idx, 1);
      }
    } catch (e) {
      CommonLog.lError(e, { at: 'openIO' });
    }
  }

  * closeIO(port: tIOPort | Array<tIOPort>): Saga<void> {
    try {
      if (port instanceof Array) {
        // eslint-disable-next-line no-restricted-syntax
        for (const p of port) {
          yield fork(ioSetApi, this.serialNumber, p.idx, 0);
        }
      } else {
        yield call(ioSetApi, this.serialNumber, port.idx, 0);
      }
    } catch (e) {
      CommonLog.lError(e, { at: 'closeIO' });
    }
  }

  * getStatus(): Saga<void> {
    try {
      yield call(ioStatusApi, this.serialNumber);
    } catch (e) {
      CommonLog.lError(e, { at: 'getStatus' });
    }
  }

  * ioContact(): Saga<void> {
    try {
      yield call(ioContactApi, this.serialNumber);
    } catch (e) {
      CommonLog.lError(e, { at: 'ioContact' });
    }
  }

  // eslint-disable-next-line flowtype/no-weak-types
  addListener(port: tIOPort, triggerMode: tIOTriggerMode, dispatcher: (...args: any) => AnyAction) {
    const listener = {
      port,
      triggerMode,
      dispatcher
    };
    this.#listeners.push(listener);
    return listener;
  }

  removeListener(listener: tIOListener) {
    return remove(this.#listeners, l => l === listener);
  }


}


