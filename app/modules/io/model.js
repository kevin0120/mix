// @flow

import { isNil } from 'lodash-es';
import Device from "../../common/type";
import type { tInputData } from "../../common/type";
import type { AnyAction } from '../../common/type';

import type { iIODataField, tIOContact, tIOData, tIOTriggerMode } from "./type";
import { CommonLog } from '../../common/utils';

/* eslint-disable no-underscore-dangle */
/* eslint-disable flowtype/no-weak-types */
export default class ClsIOModule extends Device {
  #_data: tIOData = {'inputs': {}, "outputs": {}};
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

  _storeDateField(data: tIOContact): void {
    const d = data.contact;
    if (typeof d !== 'string'){
      CommonLog.lError(`IO Data Must Be String!!!!`);
      return;
    }
    const ioType = data.type + 's';
    [...d].forEach((val, idx) => {
      this.#_data[ioType][idx].data = ClsIOModule.bitString2Boolean(val);
    })
  }

  doHandleIOData(data: tIOContact){
    const ret = this.doValidate(data.contact);
    this._storeDateField(data);

    return ret;
  }

  setIOAction(idx: number, ioType: 'inputs' | 'outputs', action: (...args: any) => AnyAction): boolean {
    if (!this._checkValidateIdx(idx, ioType)) {
      return false;
    }
    this.#_data[ioType][idx].action = action;
    return true;
  }

  set dispatcher(dispatcher: null | (...args: any) => AnyAction) {
    super.dispatcher = null; // 永远设置的是null
  }

  get dispatcher() {
    return null;
  }

  static doDispatch(data: tInputData): ?AnyAction {
    CommonLog.Info(`IO Module Please Use doIODispatch Method`);
    return null;
  }

  doIODispatch(idx: number, ioType: 'inputs' | 'outputs', ...actionParams: any): ?AnyAction {
    if (!this._checkValidateIdx(idx, ioType)) {
      return null;
    }
    const ele: iIODataField = this.#_data?.[ioType]?.[idx];
    if (isNil(ele)) {
      CommonLog.lError(`${ioType}, IO: ${idx} Is Undefined!`);
      return null;
    }
    // TODO: 触发模式引入不同的业务逻辑效果
    const triggerMode = ele.triggerMode;
    return ele.action(...actionParams);
  }

}
/* eslint-enable no-underscore-dangle */
/* eslint-enable  flowtype/no-weak-types */


