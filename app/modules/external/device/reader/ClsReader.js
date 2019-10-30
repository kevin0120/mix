// @flow

import Device from '../Device';
import { CommonLog } from '../../../../common/utils';

class ClsReader extends Device {

  // eslint-disable-next-line flowtype/no-weak-types
  constructor(...args: Array<any>) {
    super(...args);
    this.Enable();
  }

  doValidate(data: string | number): boolean {
    const ret: boolean = super.doValidate(data);
    const msg = `${this.source} validate return: ${ret.toString()}`;
    CommonLog.Info(msg);
    return ret;
  }
}

export default ClsReader;
