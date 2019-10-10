// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsScanner, { defaultScannerDispatcher } from './model';
import { CommonLog } from '../../../../common/utils';
import { getDevice } from '../index';
import type { tBarcode, tRushWebSocketData } from '../../../rush/type';

// TODO: 是否增加到系统初始化中
// scanner.Enable(); // 必须使能才能够使用

type tScannerRushData = {
  type: string,
  data: tBarcode
};

export default function* scannerNewData(data: tScannerRushData): Saga<void> {
  try {
    const d = data.data;
    CommonLog.Info(` Scanner receive data: ${d.barcode}`);
    const scanner = getDevice(d.id);
    if (scanner) {
      yield call(scanner.doDispatch, d.barcode);
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event scanner' });
  }
}
