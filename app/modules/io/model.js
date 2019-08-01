// @flow

import {isNil} from 'lodash-es';
import Device from '../../common/type'
import type {AnyAction} from '../../common/type'

import type {tIOData} from './type'
import { CommonLog } from '../../common/utils';

/* eslint-disable no-underscore-dangle */
/* eslint-disable flowtype/no-weak-types */
export default class ClsIOModule extends Device {
  _data: tIOData;

  set dispatcher(dispatcher: (...args: any) => AnyAction){
    this._dispatcher = null; // 永远设置的是null
  }

  static doDispatch(): AnyAction {
    CommonLog.Info(`IO Module Please Use doIODispatch Method`);
    return null
  }

  doIODispatch(idx: number, ioType: 'inputs' | 'outputs'): AnyAction {
      const ele = this._data?.[ioType]?.[idx];
      if (isNil(ele)){
        CommonLog.Error(`${ioType}, IO: ${idx} Is Undefined!`);
        return;
      }
      return ele.action();
  }

}
/* eslint-enable no-underscore-dangle */
/* eslint-enable  flowtype/no-weak-types */


