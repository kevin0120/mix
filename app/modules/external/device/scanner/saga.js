// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsScanner, { defaultScannerDispatcher } from './model';
import { CommonLog } from '../../../../common/utils';
import { AppendNewDevices } from '../index';
import type { tBarcode, tRushWebSocketData } from '../../../rush/type';


export const scanner = new ClsScanner('default scanner','demo scanner sn');
AppendNewDevices(scanner);
scanner.dispatcher = defaultScannerDispatcher;
// TODO: 是否增加到系统初始化中
// scanner.Enable(); // 必须使能才能够使用

type tScannerRushData={
  type:string,
  data:tBarcode
}

export default function* scannerNewData(data: tScannerRushData): Saga<void> {
  try {
    const d = data.data;
    CommonLog.Info(` Scanner receive data: ${d.barcode}`);
    yield call(scanner.doDispatch, d.barcode);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event scanner' });
  }
}


