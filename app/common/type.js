// @flow

import {Action} from 'redux'
import { CommonLog } from './utils';

const lodash = require('lodash');

type tCommonActionType = {
  +type: string
};

interface AnyAction extends Action {
  [extraProps: string]: any
}


type tDeviceNewData = {
  +data: string | number,
  +source: string
};

interface IInputDevice {
  validate(data: string | number): boolean
}

interface IOutputDevice {
  validate(data: string | number): boolean
}

// eslint-disable-next-line no-unused-vars
const defaultValidatorFunc = (data: string): boolean => true;

class Device implements IInputDevice{
  source: string;

  enable: boolean = false;

  dispatcher: () => AnyAction;

  validator: (data: string | number) => ?boolean = defaultValidatorFunc;

  constructor(source: string) {
    this.source = source;
    this.validator = null;
    this.dispatcher = null;
  }

  Enable() {
    this.enable = true;
  }

  Disable(){
    this.enable = false;
  }

  Toggle() {
    this.enable = !this.enable;
  }

  validate(data: string): boolean {
    if (lodash.isNil(data) || lodash.isEmpty(data)){
      return false
    }
    // 有效的数据
    if (lodash.isNil(this.validator)) {
      // 没有验证器默认返回正确
      return true;
    }
    return this.validator(data);
  }

  dispatch(): AnyAction {
    if (!this.enable) {
      const msg = `${this.source} Is Not Enable`;
      CommonLog(msg);
    }
    if (lodash.isNil(this.dispatcher)) {
      const msg = `${this.source} Validator is Nil, Please set Validator First`;
      CommonLog(msg);
    }
    return this.dispatcher()
  }

  setValidator(validator: (string | number) => boolean): boolean{
    this.validator = validator;
    return true;
  }

  getValidator(): ?(string | number) => boolean {
    return this.validator
  }

  setDispatcher(dispatcher: (...args: any) => AnyAction): boolean{
    this.dispatcher = dispatcher;
    return true;
  }

  removeValidator(): boolean {
    this.validator = null;
    return true
  }

  getDispatcher(): ?(...args: any) => AnyAction {
    return this.dispatcher
  }

  removeDispatcher(): boolean {
    this.dispatcher = null;
    return true
  }
}

export type {tCommonActionType,tDeviceNewData};
export default Device;
