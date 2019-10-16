import type { tDeviceNewData } from '../typeDef';
import type { tCommonActionType } from '../../../../common/type';

export const RFID = {
  READ_NEW_DATA: 'RFID_NEW_DATA'
};

export function RFIDNewData(data: string): tCommonActionType & tDeviceNewData {
  return {
    type: RFID.READ_NEW_DATA,
    data
  };
}
