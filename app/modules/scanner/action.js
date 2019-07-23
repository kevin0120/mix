// @flow

import { OPERATION_SOURCE } from '../operation/model';
import type {tCommonActionType, tDeviceNewData} from "../../common/type"

export const SCANNER = {
  READ_NEW_DATA: 'SCANNER_NEW_DATA'
};

export function ScannerNewData(aBarcode: string, source: string = OPERATION_SOURCE.SCANNER): tCommonActionType & tDeviceNewData {
  return {
    type: SCANNER.READ_NEW_DATA,
    data: aBarcode,
    source
  };
}
