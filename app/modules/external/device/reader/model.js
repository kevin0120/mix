// @flow

import Device from '../Device';
import { CommonLog } from '../../../../common/utils';
import { loginRequestUuid } from '../../../user/action';
// export const defaultReaderDispatcher = (data) => readerStepAction.getValue(data);

export type tReaderData = {
  +uid: string
};

class ClsReader extends Device {

  _dispatcher = (input) => {
    console.log(input);
    return loginRequestUuid(input.data.data.uid,'online');
  };

  constructor(...args){
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
