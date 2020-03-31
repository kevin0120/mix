// @flow

import Device from '../Device';
import { CommonLog } from '../../../common/utils';
import type {IScanner} from './IScanner';

class ClsScanner extends Device implements IScanner{
  doValidate(data: string | number): boolean {
    const ret: boolean = super.doValidate(data);
    const msg = `${this.source} validate return: ${ret.toString()}`;
    CommonLog.Info(msg);
    return ret;
  }
}

export default ClsScanner;
