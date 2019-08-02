// @flow

import Device from '../../common/type'
import { CommonLog } from '../../common/utils';

// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

class ClsReader extends Device {

  doValidate(data: string): boolean {
    const ret: boolean =  super.doValidate(data);
    const msg = `${this.source} validate return: ${ret}`;
    CommonLog.Info(msg);
    return ret
  }
}

export default ClsReader;
