// @flow

import Device from '../../common/type';
import { CommonLog } from '../../common/utils';
import { scannerStepAction } from '../step/scannerStep/action';

export const defaultScannerDispatcher = (data) => scannerStepAction.getValue(data);

class ClsScanner extends Device {

  doValidate(data: string | number ): boolean {
    const ret: boolean =  super.doValidate(data);
    const msg = `${this.source} validate return: ${ret}`;
    CommonLog.Info(msg);
    return ret;
  }
}

export default ClsScanner;
