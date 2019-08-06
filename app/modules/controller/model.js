// @flow

import Device from '../../common/type'
import { symController, AppendNewDevices } from '../global';


// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

class ClsController extends Device {

  #serialNumber: ?string = null;

  constructor(name: string, serialNumber: string){
    super(name);
    this.SerialNumber = serialNumber;
    AppendNewDevices(symController, this);
  }

  set SerialNumber(sn: string) {
    this.#serialNumber = sn
  }

  get SerialNumber() {
    return this.#serialNumber;
  }
}

export default ClsController;
