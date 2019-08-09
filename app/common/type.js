// @flow

import type { Action } from 'redux';
import { put } from 'redux-saga';
import { isURL } from 'validator/lib/isURL';
import { isNil, isEmpty, isEqual, isString } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { CommonLog } from './utils';

type tCommonActionType = {
  +type: string
};

// eslint-disable-next-line flowtype/no-weak-types
export type tInputData = any;

export type tInput = {
  data: tInputData,
  source: string,
  time: Date
};

export type AnyAction = Action & {
  // eslint-disable-next-line flowtype/no-weak-types
  [extraProps: string]: any
};


type tDeviceNewData = {
  +data: tInputData
};

// interface IInputDevice {
//   doValidate(data: string | number): boolean
// }
//
// interface IOutputDevice {
//   doValidate(data: string | number): boolean
// }

interface IHealthChecker {
  Healthz: boolean
}

class CommonExternalEntity implements IHealthChecker {
  #name: string;

  #isHealthz: boolean = false;

  #enable: boolean = false;

  constructor(name: string) {
    this.#name = name;

    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  set Healthz(isHealthz: boolean) {
    if (isEqual(this.#isHealthz, isHealthz)) {
      return;
    }
    this.#isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
    }
    const msg = `${this.#name} Healthz Status Change: ${isHealthz.toString()}`;
    CommonLog.Info(msg);
  }

  get Healthz(): boolean {
    return this.#isHealthz;
  }

  get Name(): string {
    return this.#name;
  }

  get source(): string {
    return this.#name;
  }

  get isEnable(): boolean {
    return this.#enable;
  }

  // eslint-disable-next-line require-yield
  * Enable(): Saga<void> {
    this.#enable = true;
    CommonLog.Debug(`${this.source} Is Enabled!`);
  }

  // eslint-disable-next-line require-yield
  * Disable(): Saga<void> {
    this.#enable = false;
    CommonLog.Debug(`${this.source} Is Disabled!`);
  }

  // eslint-disable-next-line require-yield
  * ToggleEnable(): Saga<void> {
    this.#enable = !this.#enable;
  }

}

class ExternalSystem extends CommonExternalEntity {
  #endpoint: ?string = null;

  constructor(name: string, endpoint: string) {
    super(name);
    this.Endpoint = endpoint;
  }

  get Endpoint(): ?string {
    return this.#endpoint;
  }

  set Endpoint(endpoint: string) {
    if (!isURL(endpoint)) {
      CommonLog.lError(`Endpoint: ${endpoint} Is Not Valid!`);
    } else {
      this.#endpoint = endpoint;
    }
  }

  Connect(): boolean {
    CommonLog.Warn('Please Override Connect Method!!!');
    this.Healthz = false;
    return false;
  }

  Close(): boolean {
    CommonLog.Warn('Please Override Close Method!!!');
    this.Healthz = false;
    return true;
  }

}

// eslint-disable-next-line no-unused-vars
const defaultValidatorFunc = (data: string | number): boolean => true;

/* eslint-disable no-underscore-dangle */

class Device extends CommonExternalEntity {

  #dispatcher: null | (?tInput) => AnyAction = null;

  #validator: null | (data: tInputData) => boolean = defaultValidatorFunc;

  // eslint-disable-next-line flowtype/no-weak-types
  constructor(name: string) {
    super(name);
    // eslint-disable-next-line flowtype/no-weak-types
    (this: any).doDispatch = this.doDispatch.bind(this);
  }

  doValidate(data: tInputData): boolean {
    let _isEmpty = false;
    if (isNil(data)) {
      _isEmpty = true;
    }
    if (isString(data) && isEmpty(data)) {
      _isEmpty = true;
    }
    if (_isEmpty) {
      const msg = `${this.source} Receive Empty Data: ${data}`;
      CommonLog.Debug(msg);
      return false;
    }
    // 有效的数据
    if (!this.validator) {
      // 没有验证器默认返回正确
      return true;
    }
    return this.validator(data);
  }

  * doDispatch(data: tInputData): Saga<void> {
    try {
      if (!this.isEnable) {
        const msg = `${this.source} Is Not Enabled`;
        CommonLog.Info(msg);
        return;
      }
      if(!this.doValidate(data)){
        const msg = `data ${data} didn't pass validator at ${this.source}`;
        CommonLog.Info(msg);
        return ;
      }
      if (!this.dispatcher) {
        const msg = `${this.source} Validator is Nil, Please set Validator First`;
        CommonLog.Warn(msg);
        return;
      }
      yield put(this.dispatcher({
        data,
        source: this.Name,
        time: new Date()
      }));
    } catch (e) {
      CommonLog.lError(e, {
        at: 'doDispatch',
        data
      });
    }
  }

  set validator(validator: (string | number) => boolean) {
    this.#validator = validator;
  }

  get validator(): ?(string | number) => boolean {
    return this.#validator;
  }

  // eslint-disable-next-line flowtype/no-weak-types
  set dispatcher(dispatcher: null | (?tInput) => AnyAction) {
    this.#dispatcher = dispatcher;
  }

  // eslint-disable-next-line flowtype/no-weak-types
  get dispatcher(): ?(?tInput) => AnyAction {
    return this.#dispatcher;
  }

  RemoveValidator(): boolean {
    this.#validator = null;
    return true;
  }

  RemoveDispatcher(): boolean {
    this.#dispatcher = null;
    return true;
  }
}
/* eslint-enable no-underscore-dangle */

export type { tCommonActionType, tDeviceNewData };
export default Device;
