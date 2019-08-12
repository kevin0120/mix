// @flow

import { call, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsScanner, { defaultScannerDispatcher } from './model';
import { CommonLog } from '../../../../common/utils';
// eslint-disable-next-line import/named
import { symScanner, AppendNewDevices } from '../global';
import type { tBarcode, tRushWebSocketData } from '../../../rush/type';


export const scanner = new ClsScanner(symScanner);
AppendNewDevices(symScanner, scanner);
scanner.dispatcher = defaultScannerDispatcher;
// TODO: 是否增加到系统初始化中
// scanner.Enable(); // 必须使能才能够使用

export default function* scannerNewData(data: tRushWebSocketData): Saga<void> {
  try {
    const d = (data.data: tBarcode);
    CommonLog.Info(` Scanner receive data: ${d.barcode}`);
    yield call(scanner.doDispatch, d.barcode);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event scanner' });
  }
}


