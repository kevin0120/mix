// @flow

import Device from '../../common/type'
import { symTool, AppendNewDevices } from '../global';
import ClsController from '../controller/model';
import type { tInputData } from '../../common/type';
import screwStepAction from '../step/screwStep/action';
import { CommonLog } from '../../common/utils';

// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

export const defaultScrewToolDispatcher = (data: tInputData) => screwStepAction.result(data);

class ClsScrewTool extends Device {

  #serialNumber: ?string = null;

  #belongTo: ?ClsController = null;

  constructor(name: string, serialNumber: string){
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
    this.#belongTo = controller
  }

  get BelongTo() {
    return this.#belongTo
  }

  set SerialNumber(sn: string) {
    this.#serialNumber = sn
  }

  get SerialNumber() {
    return this.#serialNumber;
  }
}

export default ClsScrewTool;


