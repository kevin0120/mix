// @flow

import {Action} from 'redux'
import { CommonLog } from './utils';

const lodash = require('lodash');

type tCommonActionType = {
  +type: string
};

/* eslint-disable flowtype/no-weak-types */
export interface AnyAction extends Action {
  [extraProps: string]: any
}

type tDeviceNewData = {
  +data: {[key: string]: any}
};
/* eslint-enable flowtype/no-weak-types */

interface IInputDevice {
  validate(data: string | number): boolean
}

interface IOutputDevice {
  validate(data: string | number): boolean
}

// eslint-disable-next-line no-unused-vars
const defaultValidatorFunc = (data: string): boolean => true;

/* eslint-disable no-underscore-dangle */

class Device implements IInputDevice{
  _source: string;

  _isHealthz: boolean = false;

  _enable: boolean = false;

  _dispatcher: () => AnyAction;

  _validator: (data: string | number) => ?boolean = defaultValidatorFunc;

  constructor(source: string) {
    this._source = source;
    this._validator = null;
    this._dispatcher = null;
  }

  get source(): string {
    return this._source;
  }

  set Healthz(isHealthz: boolean) {
    if (lodash.isEqual(this._isHealthz,isHealthz)){
      return
    }
    this._isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
    }
    const msg = `${this._source} Healthz Status Change: ${isHealthz}`;
    CommonLog.Info(msg)
  }

  Enable() {
    this._enable = true;
  }

  Disable(){
    this._enable = false;
  }

  ToggleEnable() {
    this._enable = !this._enable;
  }

  doValidate(data: string): boolean {
    if (lodash.isNil(data) || lodash.isEmpty(data)){
      const msg = `${this._source} Receive Empty Data: ${data}`;
      CommonLog.Debug(msg);
      return false
    }
    // 有效的数据
    if (lodash.isNil(this.validator)) {
      // 没有验证器默认返回正确
      return true;
    }
    return this.validator(data);
  }

  doDispatch(data: string): AnyAction {
    if (!this._enable) {
      const msg = `${this._source} Is Not Enable`;
      CommonLog.Info(msg);
    }
    if (lodash.isNil(this.dispatcher)) {
      const msg = `${this._source} Validator is Nil, Please set Validator First`;
      CommonLog.Warn(msg);
    }
    return this.dispatcher(data);
  }

  set validator(validator: (string | number) => boolean){
    this._validator = validator;
  }

  get validator(): ?(string | number) => boolean {
    return this._validator
  }

  /* eslint-disable flowtype/no-weak-types */
  set dispatcher(dispatcher: (...args: any) => AnyAction){
    this._dispatcher = dispatcher;
  }

  get dispatcher(): ?(...args: any) => AnyAction {
    return this._dispatcher
  }
  /* eslint-enable flowtype/no-weak-types */

  RemoveValidator(): boolean {
    this._validator = null;
    return true
  }

  RemoveDispatcher(): boolean {
    this._dispatcher = null;
    return true
  }
}

/* eslint-enable no-underscore-dangle */

export type {tCommonActionType,tDeviceNewData};
export default Device;
