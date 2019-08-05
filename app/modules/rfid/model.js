// @flow

import Device from '../../common/type';

class ClsRFID extends Device {
  validate(data: string): boolean {
    return super.doValidate(data);
  }
}

export default ClsRFID;
