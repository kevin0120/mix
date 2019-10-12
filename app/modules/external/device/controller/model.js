// @flow

import Device from '../Device';

class ClsController extends Device {
  #serialNumber: ?string = null;

  constructor(name: string, serialNumber: string, ...rest: Array<any>) {
    super(name, serialNumber);
    this.#serialNumber = serialNumber;
  }

  set serialNumber(sn: string) {
    this.#serialNumber = sn;
  }

  get serialNumber() {
    return this.#serialNumber;
  }
}

export default ClsController;
