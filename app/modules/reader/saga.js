// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import {READER} from './action';
import ClsReader  from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';
import { gDevices } from '../global';

const symReader = 'Reader';

export const reader = new ClsReader(symReader);
gDevices[symReader] = [reader];
// reader.dispatcher = defaultReaderDispatcher;

function* readerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (reader.doValidate(data)){
      yield put(reader.doDispatch(data))
    } else {
      // do nothing
    }
  } catch (e) {
    CommonLog.Error(e)
  }
}

export default function* watchReader(): Saga<void> {
  yield takeLatest(READER.READ_NEW_DATA, readerHandler);
}
