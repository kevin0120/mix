// @flow

import { isNil } from 'lodash-es';
import Device from '../../common/type'
import { gDevices } from '../global';
import ClsController from '../controller/model';

const symTool = 'ToolScrew';

// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

class ClsScrewTool extends Device {

  #serialNumber: ?string = null;

  #belongTo: ?ClsController = null;

  constructor(name: string, serialNumber: string){
    super(name);
    this.SerialNumber = serialNumber;
    const cl = gDevices?.[symTool];
    if (isNil(cl)){
      gDevices[symTool] = [this];
    }else {
      cl.push(this);
    }
  }

  set SerialNumber(sn: string) {
    this.#serialNumber = sn
  }

  get SerialNumber() {
    return this.#serialNumber;
  }
}

export default ClsScrewTool;


