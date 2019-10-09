// @flow

import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import Device from '../Device';
// import ClsController from '../controller/model';
import type { AnyAction, tInputData } from '../type';
import screwStepAction from '../../../step/screwStep/action';
import { CommonLog } from '../../../../common/utils';
import { toolEnableApi } from '../../../../api/tools';

export const defaultScrewToolDispatcher = (data: tInputData): AnyAction =>
  screwStepAction.result(data);

export default class ClsScrewTool extends Device {
  constructor(name: string, serialNumber: string) {
    super(name, serialNumber);
    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    this.dispatcher = defaultScrewToolDispatcher;
    /* eslint-enable flowtype/no-weak-types */
  }

  doValidate(data: string | number): boolean {
    const ret: boolean = super.doValidate(data);
    const msg = `${this.source} validate return: ${ret.toString()}`;
    CommonLog.Info(msg);
    return ret;
  }

  *Enable(): Saga<void> {
    try {
      if (!this.isEnable) {
        const { result, msg } = yield call(
          toolEnableApi,
          this.serialNumber,
          true
        );
        if (result !== 0) {
          CommonLog.lError(`tool enable failed:${msg}`, {
            at: 'ClsScrewTool.Enable',
            tool_sn: this.serialNumber,
            enable: true
          });
          return false;
        }
        yield call([this, super.Enable]);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Enable'
      });
    }
  }

  *Disable(): Saga<void> {
    try {
      if (this.isEnable) {
        const { result, msg } = yield call(
          toolEnableApi,
          this.serialNumber,
          false
        );
        if (result !== 0) {
          CommonLog.lError(`tool enable failed:${msg}`, {
            at: 'ClsScrewTool.Disable',
            tool_sn: this.serialNumber,
            enable: false
          });
          return false;
        }
        yield call([this, super.Disable]);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Disable'
      });
    }
  }

  *ToggleEnable(): Saga<void> {
    try {
      const { result, msg } = yield call(
        toolEnableApi,
        this.serialNumber,
        !this.isEnable
      );
      if (result !== 0) {
        CommonLog.lError(`tool enable failed:${msg}`, {
          at: 'ClsScrewTool.ToggleEnable',
          tool_sn: this.serialNumber,
          enable: !this.isEnable
        });
        return false;
      }
      yield call([this, super.ToggleEnable]);
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.ToggleEnable'
      });
    }
  }
}
