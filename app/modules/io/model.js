// @flow

import { isNil } from 'lodash-es';
import Device from '../../common/type';
import type { AnyAction } from '../../common/type';

import type { tIOData } from './type';
import { CommonLog } from '../../common/utils';

/* eslint-disable no-underscore-dangle */
/* eslint-disable flowtype/no-weak-types */
export default class ClsIOModule extends Device {
  #_data: tIOData;
  #_maxInputs: number = 0;
  #_maxOutputs: number = 0;

  constructor(name: string, inputs: number, outputs: number) {
    super(name);
    this.#_maxInputs = inputs;
    this.#_maxOutputs = outputs;

  }

  _checkValidateIdx(idx: number, ioType: 'inputs' | 'outputs'): boolean {
    switch (ioType) {
      case 'inputs':
        return idx <= this.#_maxInputs;
      case 'outputs':
        return idx <= this.#_maxOutputs;
      default:
        return false;
    }
  }

  setIOAction(idx: number, ioType: 'inputs' | 'outputs', action: (...args: any) => AnyAction): boolean {
    if (!this._checkValidateIdx(idx, ioType)) {
      return false;
    }
    this.#_data[ioType][idx].action = action;
    return true;
  }

  set dispatcher(dispatcher: (...args: any) => AnyAction) {
    super.dispatcher = null; // 永远设置的是null
  }

  static doDispatch(): AnyAction {
    CommonLog.Info(`IO Module Please Use doIODispatch Method`);
    return null;
  }

  doIODispatch(idx: number, ioType: 'inputs' | 'outputs'): ?AnyAction {
    if (!this._checkValidateIdx(idx, ioType)) {
      return null;
    }
    const ele = this.#_data?.[ioType]?.[idx];
    if (isNil(ele)) {
      CommonLog.lError(`${ioType}, IO: ${idx} Is Undefined!`);
      return null;
    }
    return ele.action();
  }

}
/* eslint-enable no-underscore-dangle */
/* eslint-enable  flowtype/no-weak-types */


