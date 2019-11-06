// @flow

import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import Device from '../Device';
import { CommonLog } from '../../../common/utils';
import { toolEnableApi } from '../../../api/tools';
import type { IScrewTool } from './interface/IScrewTool';

export default class ClsScrewTool extends Device implements IScrewTool {
  constructor(name: string, serialNumber: string, ...rest: Array<any>) {
    super(name, serialNumber);
    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  doValidate(data: string | number): boolean {
    const ret: boolean = super.doValidate(data);
    const msg = `${this.source} validate return: ${ret.toString()}`;
    CommonLog.Info(msg);
    return ret;
  }

  * Enable(): Saga<void> {
    try {
      if (!this.isEnable) {
        yield call(
          (toolEnableApi: Function),
          this.serialNumber,
          true
        );
        yield call([this, super.Enable]);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Enable'
      });
      throw new Error(`工具使能失败: ${e.message}`);

    }
  }

  * Disable(): Saga<void> {
    try {
      if (this.isEnable) {
        yield call(
          toolEnableApi,
          this.serialNumber || '',
          false
        );
        yield call([this, super.Disable]);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Disable'
      });
      throw new Error(`工具禁用失败: ${e.message}`);
    }
  }

  * ToggleEnable(): Saga<void> {
    try {
      yield call(
        toolEnableApi,
        this.serialNumber || '',
        !this.isEnable
      );
      yield call([this, super.ToggleEnable]);
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.ToggleEnable'
      });
      throw new Error(`切换工具使能失败: ${e.message}`);
    }
  }
}
