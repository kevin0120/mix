// @flow
import type { tDeviceNewData } from '../typeDef';
import type { tCommonActionType } from '../../../common/type';

export const RFID = {
  READ_NEW_DATA: 'RFID_NEW_DATA',
  TOGGLE_STATUS: 'RFID_TOGGLE_STATUS'
};

export function RFIDNewData(data: string): tCommonActionType & tDeviceNewData {
  return {
    type: RFID.READ_NEW_DATA,
    data
  };
}
