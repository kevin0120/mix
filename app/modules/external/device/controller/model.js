// @flow

import Device from '../Device';

class ClsController extends Device {
  _serialNumber: ?string = null;

  // eslint-disable-next-line no-unused-vars,flowtype/no-weak-types
  constructor(name: string, serialNumber: string, ...rest: Array<any>) {
    super(name, serialNumber);
    this._serialNumber = serialNumber;
  }

  set serialNumber(sn: string) {
    this._serialNumber = sn;
  }

  get serialNumber() {
    return this._serialNumber;
  }
}

export default ClsController;
