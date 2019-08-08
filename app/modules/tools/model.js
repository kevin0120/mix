// @flow

import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import Device from '../../common/type';
import { symTool, AppendNewDevices } from '../global';
import ClsController from '../controller/model';
import type { tInputData } from '../../common/type';
import screwStepAction from '../step/screwStep/action';
import { CommonLog } from '../../common/utils';
import { toolEnableApi } from '../../api/order';

export const defaultScrewToolDispatcher = (data: tInputData) => screwStepAction.result(data);

class ClsScrewTool extends Device {

  #serialNumber: ?string = null;

  #belongTo: ?ClsController = null;

  constructor(name: string, serialNumber: string) {
    super(name);
    this.SerialNumber = serialNumber;
    AppendNewDevices(symTool, this);
  }

  doValidate(data: string | number): boolean {
    const ret: boolean = super.doValidate(data);
    const msg = `${this.source} validate return: ${ret.toString()}`;
    CommonLog.Info(msg);
    return ret;
  }

  set BelongTo(controller: ClsController) {
    this.#belongTo = controller;
  }

  get BelongTo() {
    return this.#belongTo;
  }

  set SerialNumber(sn: string) {
    this.#serialNumber = sn;
  }

  get SerialNumber() {
    return this.#serialNumber;
  }

  * Enable(): Saga<void> {
    try {
      if (!this.isEnable) {
        const { result, msg } = yield call(toolEnableApi, this.SerialNumber, true);
        if (result !== 0) {
          CommonLog.lError(`tool enable failed:${msg}`, {
            at: 'ClsScrewTool.Enable',
            tool_sn: this.SerialNumber,
            enable: true
          });
          return false;
        }
        yield call(super.Enable);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Enable'
      });
    }
  }

  * Disable(): Saga<void> {
    try {
      if (this.isEnable) {
        const { result, msg } = yield call(toolEnableApi, this.SerialNumber, false);
        if (result !== 0) {
          CommonLog.lError(`tool enable failed:${msg}`, {
            at: 'ClsScrewTool.Disable',
            tool_sn: this.SerialNumber,
            enable: false
          });
          return false;
        }
        yield call(super.Disable);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.Disable'
      });
    }
  }

  * ToggleEnable(): Saga<void> {
    try {
      const { result, msg } = yield call(toolEnableApi, this.SerialNumber, !this.isEnable);
      if (result !== 0) {
        CommonLog.lError(`tool enable failed:${msg}`, {
          at: 'ClsScrewTool.ToggleEnable',
          tool_sn: this.SerialNumber,
          enable: !this.isEnable
        });
        return false;
      }
      yield call(super.ToggleEnable);
    } catch (e) {
      CommonLog.lError(e, {
        at: 'ClsScrewTool.ToggleEnable'
      });
    }
  }
}

export default ClsScrewTool;


