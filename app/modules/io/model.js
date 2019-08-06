// @flow

import { isNil } from 'lodash-es';
import Device from "../../common/type";
import type { tInputData, AnyAction } from "../../common/type";

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
    const {type: dType, contact} = data;
    if (typeof contact !== 'string'){
      CommonLog.lError(`IO Data Must Be String!!!!`);
      return;
    }
    const ioType = dType + 's';
    [...contact].forEach((val, idx) => {
      this.#_data[ioType][idx].data = ClsIOModule.bitString2Boolean(val);
    })
  }

  static _getTriggerMode(old: boolean, last: boolean, mode: tIOTriggerMode): ?tIOTriggerMode {
    let ret = mode;
    if(old === last){
      return null;
    }
    return ret;
  }

  _doHandleIODataPart1(data: tIOContact,...actionParams: any): boolean {
    let ret = false;
    const {type: dType, contact} = data;
    const ioType = ((dType + 's': any): 'inputs' | 'outputs');
    const eles = this.#_data?.[ioType];
    if (isNil(eles)){
      return ret;
    }
    [...contact].map((val, idx) => {
      if(isNil(val)){
        return
      }
      if (!this._checkValidateIdx(idx, ioType)) {
        return null;
      }
      const e: iIODataField = this.#_data[ioType]?.[idx];
      if(isNil(e)){
        // TODO: 未定义以默认方式进行
      }else {
        const old = e.data;
        const last = ClsIOModule.bitString2Boolean(val);
        const mode = ClsIOModule._getTriggerMode(old, last, e.triggerMode);
        if(!isNil(mode)){
          this._doIODispatch(e, ...actionParams);
        }
      }
    });

    return ret;

  }

  doHandleIOData(data: tIOContact, ...actionParams: any){
    const ret = this.doValidate(data.contact);
    this._doHandleIODataPart1(data, ...actionParams);
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
    CommonLog.Info(`IO Module Please Use doHandleIOData Method`);
    return null;
  }

  _doIODispatch(ele: iIODataField, ...actionParams: any): ?AnyAction {
    return ele.action(...actionParams);
  }

}
/* eslint-enable no-underscore-dangle */
/* eslint-enable  flowtype/no-weak-types */


