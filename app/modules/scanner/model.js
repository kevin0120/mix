// @flow

import {Device, IInputDevice } from '../../common/type'

const lodash = require('lodash');

class Scanner extends Device implements IInputDevice {
  static validate(data: string): boolean {
    if (lodash.isNil(data) || lodash.isEmpty(data)){
      return false
    }
    // 有效的数据
    return true
  }
}

export default Scanner;
