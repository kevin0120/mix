// @flow

import type {AnyAction} from 'redux'

const lodash = require('lodash');

type tCommonActionType = {
  +type: string
};


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

class Device implements IInputDevice{
  source: string;

  dispatcher: () => AnyAction;

  validator: (data: string | number) => ?boolean = null;

  constructor(source: string) {
    this.source = source;
    this.validator = null;
    this.dispatcher = null;
  }

  validate(data: string): boolean {
    if (lodash.isNil(data) || lodash.isEmpty(data)){
      return false
    }
    // 有效的数据
    if (lodash.isNil(this.validator)) {
      console.error("Scanner Validator is Nil, Please set Validator First")
      return false;
    }
    return this.validator(data);
  }

  dispatch(): AnyAction {
    if (lodash.isNil(this.dispatcher)) {
      console.error("Scanner Validator is Nil, Please set Validator First")
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
