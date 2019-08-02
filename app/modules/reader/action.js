// @flow

import type {tCommonActionType, tDeviceNewData} from "../../common/type"

export const READER = {
  READ_NEW_DATA: 'READER_NEW_DATA'
};

export function ReaderNewData(uid: string): tCommonActionType & tDeviceNewData {
  return {
    type: READER.READ_NEW_DATA,
    data: uid,
  };
}
