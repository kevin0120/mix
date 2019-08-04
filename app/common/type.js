// @flow

import { Action } from 'redux';
import { CommonLog } from './utils';
import { isURL } from 'validator/lib/isURL';
import { isNil, isEmpty, isEqual, isString } from 'lodash-es';

type tCommonActionType = {
  +type: string
};

/* eslint-disable flowtype/no-weak-types */
export interface AnyAction extends Action {
  [extraProps: string]: any
}

type tDeviceNewData = {
  +data: { [key: string]: any }
};
/* eslint-enable flowtype/no-weak-types */

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
  }

  set Healthz(isHealthz: boolean) {
    if (isEqual(this.#isHealthz, isHealthz)) {
      return;
    }
    this.#isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
    }
    const msg = `${this.#name} Healthz Status Change: ${isHealthz}`;
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

  Enable() {
    CommonLog.Debug(`${this.source} Is Enable!`);
    this.#enable = true;
  }

  Disable() {
    CommonLog.Debug(`${this.source} Is Disable!`);
    this.#enable = false;
  }

  ToggleEnable() {
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
    return false
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
  #dispatcher: null | () => AnyAction = null;

  #validator: null | (data: string | number) => ?boolean = defaultValidatorFunc;

  doValidate(data: string | number): boolean {
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
    if (isNil(this.validator)) {
      // 没有验证器默认返回正确
      return true;
    }
    return this.validator(data);
  }

  doDispatch(data: string): AnyAction {
    if (!this.#enable) {
      const msg = `${this.source} Is Not Enable`;
      CommonLog.Info(msg);
      return
    }
    if (isNil(this.dispatcher)) {
      const msg = `${this.source} Validator is Nil, Please set Validator First`;
      CommonLog.Warn(msg);
      return
    }
    return this.dispatcher(data);
  }

  set validator(validator: (string | number) => boolean) {
    this.#validator = validator;
  }

  get validator(): ?(string | number) => boolean {
    return this.#validator;
  }

  /* eslint-disable flowtype/no-weak-types */
  set dispatcher(dispatcher: (...args: any) => AnyAction) {
    this.#dispatcher = dispatcher;
  }

  get dispatcher(): ?(...args: any) => AnyAction {
    return this.#dispatcher;
  }

  /* eslint-enable flowtype/no-weak-types */

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
