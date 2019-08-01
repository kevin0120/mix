// @flow

import type {tCommonActionType, tDeviceNewData} from "../../common/type"

export const SCANNER = {
  READ_NEW_DATA: 'SCANNER_NEW_DATA'
};

export function ScannerNewData(aBarcode: string): tCommonActionType & tDeviceNewData {
  return {
    type: SCANNER.READ_NEW_DATA,
    data: aBarcode,
  };
}
