// @flow

import { isNil, remove } from 'lodash-es';
import Device from '../../common/type';
import type { AnyAction } from '../../common/type';
import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import type {
  tIOContact,
  tIOData,
  tIOTriggerMode,
  tIODirection,
  tIOListener,
  tIOPort,
  tIOChange
} from './type';
import { CommonLog } from '../../common/utils';
import { AppendNewDevices, symIO } from '../global';
import { ioDirection, ioTriggerMode } from './type';
import { ioSetApi, ioContactApi, ioStatusApi } from '../../api/io';

export default class ClsIOModule extends Device {
  #_data: tIOData = { input: '', output: '' };
  #_ports: Array<tIOPort> = [];
  #_maxInputs: number = 0;
  #_maxOutputs: number = 0;
  #_listeners: Array<tIOListener> = [];
  #serialNumber: ?string = null;

  constructor(name: string, serialNumber: string, inputs: number, outputs: number) {
    super(name);
    this.#serialNumber = serialNumber;
    this.#_maxInputs = inputs;
    this.#_maxOutputs = outputs;
    for (let i = 0; i < inputs; i++) {
      this.#_ports.push({
        direction: ioDirection.input,
        idx: i
      });
    }
    for (let i = 0; i < outputs; i++) {
      this.#_ports.push({
        direction: ioDirection.output,
        idx: i
      });
    }
    AppendNewDevices(symIO, this);
  }

  getPort(direction: tIODirection, idx: number): ?tIOPort {
    return this.#_ports.find(p => p.direction === direction && p.idx === idx);
  }

  _checkValidateIdx(idx: number, ioType: tIODirection): boolean {
    switch (ioType) {
      case 'input':
        return idx <= this.#_maxInputs;
      case 'output':
        return idx <= this.#_maxOutputs;
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
    // [...contact].forEach((val, idx) => {
    //   this.#_data[direction][idx].data = ClsIOModule.bitString2Boolean(val);
    // });
    this.#_data[direction] = contact;
  }

  // static _getTriggerMode(old: boolean, last: boolean, mode: tIOTriggerMode): ?tIOTriggerMode {
  //   let ret = mode;
  //   if (old === last) {
  //     return null;
  //   }
  //   return ret;
  // }
  get serialNumber() {
    return this.#serialNumber;
  }

  _getIOChanges(newData: tIOContact): Array<tIOChange> {
    const changes: Array<tIOChange> = [];
    const { direction, contact } = newData;
    const eles: string = this.#_data?.[direction];
    if (!eles) {
      return changes;
    }
    const startIdx = direction === ioDirection.input ? 0 : this.#_maxInputs;

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
          port: this.#_ports[startIdx + idx],
          triggerMode: ioTriggerMode.change
        });
        if (newBool) {
          changes.push({
            port: this.#_ports[startIdx + idx],
            triggerMode: ioTriggerMode.rising
          });
        } else {
          changes.push({
            port: this.#_ports[startIdx + idx],
            triggerMode: ioTriggerMode.falling
          });
        }
      }

      // const e: iIODataField = this.#_data[direction]?.[idx];
      // if (isNil(e)) {
      //   // TODO: 未定义以默认方式进行
      // } else {
      //   const old = e.data;
      //   const last = ClsIOModule.bitString2Boolean(val);
      //   const mode = ClsIOModule._getTriggerMode(old, last, e.triggerMode);
      //   if (!isNil(mode)) {
      //     this._doIODispatch(e, ...actionParams);
      //   }
      // }
    });

    return changes;

  }

  _doHandleIOData(newData: tIOContact, ...actionParams: any): Array<AnyAction> {
    const ret = this.doValidate(newData.contact);
    if (!ret) {
      return [];
    }
    const changes = this._getIOChanges(newData, ...actionParams);
    const matchedListeners = this.#_listeners.filter((l) =>
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

  // setIOAction(idx: number, ioType: tIODirection, action: (...args: any) => AnyAction): boolean {
  //   if (!this._checkValidateIdx(idx, ioType)) {
  //     return false;
  //   }
  //   this.#_data[ioType][idx].action = action;
  //   return true;
  // }

  set dispatcher(dispatcher: null | (...args: any) => AnyAction) {
    super.dispatcher = null; // 永远设置的是null
  }

  get dispatcher() {
    return null;
  }

  doDispatch(newData: tIOContact): Array<AnyAction> {
    // CommonLog.Info(`IO Module Please Use doHandleIOData Method`);
    return this._doHandleIOData(newData);
  }

  // _doIODispatch(ele: iIODataField, ...actionParams: any): ?AnyAction {
  //   return ele.action(...actionParams);
  // }

  * setIO(port: tIOPort, value: boolean): Saga<void> {
    // TODO: handle set io
    const status = value ? 1 : 0;
    yield call(ioSetApi, this.serialNumber, port.idx, status);
  }

  * getStatus(): Saga<void> {
    yield call(ioStatusApi, this.serialNumber);
  }

  * ioContact(): Saga<void> {
    yield call(ioContactApi, this.serialNumber);
  }

  addListener(port: tIOPort, triggerMode: tIOTriggerMode, dispatcher: (...args: any) => AnyAction) {
    const listener = {
      port,
      triggerMode,
      dispatcher
    };
    this.#_listeners.push(listener);
    return listener;
  }

  removeListener(listener: tIOListener) {
    return remove(this.#_listeners, l => l === listener);
  }


}


