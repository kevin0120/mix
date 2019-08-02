// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { SCANNER } from './action';
import ClsScanner from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';

// eslint-disable-next-line prefer-const
let scanner = new ClsScanner('clsScanner');

function* scannerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (scanner.validate(data)) {
      yield put(scanner.doDispatch());
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
