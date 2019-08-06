// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { SCANNER } from './action';
import ClsScanner, { defaultScannerDispatcher } from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';
// eslint-disable-next-line import/named
import { symScanner, AppendNewDevices } from '../global';


export const scanner = new ClsScanner(symScanner);
AppendNewDevices(symScanner, scanner);
scanner.dispatcher = defaultScannerDispatcher;
// TODO: 是否增加到系统初始化中
// scanner.Enable(); // 必须使能才能够使用

function* scannerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (scanner.doValidate(data)) {
      const respAction = scanner.doDispatch(data);
      if (respAction) {
        yield put(respAction);
      }
    } else {
      // do nothing
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

export default function* watchScanner(): Saga<void> {
  yield takeLatest(SCANNER.READ_NEW_DATA, scannerHandler);
}
