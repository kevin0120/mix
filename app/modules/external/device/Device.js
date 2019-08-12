import { isEmpty, isNil, isString } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { put } from 'redux-saga/effects';
import { CommonLog } from '../../../common/utils';
import type { AnyAction, tInput, tInputData } from './type';
import CommonExternalEntity from '../CommonExternalEntity';

const defaultValidatorFunc = (data: string | number): boolean => true;

export default class Device extends CommonExternalEntity{

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