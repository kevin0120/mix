// @flow

import Device from '../Device';

class ClsRFID extends Device {
  validate(data: string): boolean {
    return super.doValidate(data);
  }
}

export default ClsRFID;
