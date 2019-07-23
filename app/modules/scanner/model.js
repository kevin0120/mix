// @flow

import { Device } from '../../common/type'
import { Info } from '../../logger';

class Scanner extends Device {
  validate(data: string): boolean {
    const ret: boolean =  super.validate(data);
    Info(`Scanner validate return: ${ret}`);
    return ret
  }

}

export default Scanner;
