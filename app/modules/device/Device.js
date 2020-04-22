// @flow
import { isEmpty, isNil, isString } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { put, all, select } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import type {
  tInput,
  tInputData,
  tDeviceSN,
  tInputPredicate,
  tInputListener
} from './typeDef';
import CommonExternalEntity from '../external/CommonExternalEntity';
import type { IDevice } from './IDevice';
import type { tAction, tListenerObj } from '../typeDef';
import { makeListener } from '../util';

const defaultValidatorFunc = (): boolean => true;

export default class Device extends CommonExternalEntity implements IDevice {
  _validator: null | ((data: tInputData) => boolean) = defaultValidatorFunc;

  _serialNumber: ?tDeviceSN = null;

  _inputListener: tListenerObj = makeListener();

  // eslint-disable-next-line flowtype/no-weak-types,no-unused-vars
  constructor(name: string, sn: tDeviceSN, config: Object, data: any) {
    super(name, config);
    this._serialNumber = sn;
    // eslint-disable-next-line flowtype/no-weak-types
    (this: any).doDispatch = this.doDispatch.bind(this);
  }

  // eslint-disable-next-line flowtype/no-weak-types
  addListener(
    predicate: tInputPredicate,
    action: (...args: any) => tAction<any, any>
  ): tInputListener {
    return this._inputListener.add(predicate, action);
  }

  removeListener(listener: tInputListener): Array<tInputListener> {
    this._inputListener.remove(listener);
    // console.log(this._inputListener.listeners);
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
      CommonLog.Info(msg);
      return false;
    }
    // 有效的数据
    if (!this.validator) {
      // 没有验证器默认返回正确
      return true;
    }
    return this.validator(data);
  }

  *doDispatch(data: tInputData): Saga<void> {
    try {
      console.log(data);
      // if (!this.isEnable) {
      //   const msg = `${this.source} Is Not Enabled`;
      //   CommonLog.Info(msg);
      //   return;
      // }
      if (!this.doValidate(data)) {
        const msg = `data ${data} didn't pass validator at ${this.source}`;
        CommonLog.Info(msg);
        return;
      }
      if (!this._inputListener) {
        const msg = `${this.source} listener is Nil, Please set listener First`;
        CommonLog.Warn(msg);
        return;
      }

      const input: tInput = {
        data,
        source: this.source,
        time: new Date()
      };
      const state = yield select();
      const actions = this._inputListener.check(input, state);

      yield all(actions.map(a => put(a)));
    } catch (e) {
      CommonLog.lError(e, {
        at: 'doDispatch',
        data
      });
      throw e;
    }
  }

  set validator(validator: (string | number) => boolean) {
    this._validator = validator;
  }

  get validator(): ?(string | number) => boolean {
    return this._validator;
  }

  RemoveValidator(): boolean {
    this._validator = null;
    return true;
  }

  get serialNumber() {
    return this._serialNumber;
  }
}
