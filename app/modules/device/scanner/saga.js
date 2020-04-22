// @flow

import { call ,put} from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../../common/utils';
import { getDevice } from '../../deviceManager/devices';
import type { tBarcode } from '../../rush/type';
import  { newData } from '../../manual/action';


type tScannerRushData = {
  type: string,
  data: tBarcode
};

export default function* scannerNewData(data: tScannerRushData): Saga<void> {
  try {
    const d = data.data;
    yield put(newData(d.barcode));
    CommonLog.Info(` Scanner receive data: ${d.barcode}`);
    const scanner = getDevice(d.sn);
    if (scanner) {
      yield call(scanner.doDispatch, d.barcode);
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event scanner' });
  }
}
