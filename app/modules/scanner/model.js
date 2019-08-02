// @flow

import Device from '../../common/type';
import { CommonLog } from '../../common/utils';

class ClsScanner extends Device {
  validate(data: string): boolean {
    const ret: boolean = super.validate(data);
    const msg = `${this.source} validate return: ${ret}`;
    CommonLog.Info(msg);
    return ret;
  }
}

export default ClsScanner;
