// @flow

import { isNil } from 'lodash-es';
import Device from '../../common/type'
import { gDevices } from '../global';

const symController = 'Controller';

// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

class ClsController extends Device {

  #serialNumber: ?string = null;

  constructor(name: string, serialNumber: string){
    super(name);
    this.SerialNumber = serialNumber;
    const cl = gDevices?.[symController];
    if (isNil(cl)){
      gDevices[symController] = [this];
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

export default ClsController;
