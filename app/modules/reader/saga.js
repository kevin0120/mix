// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import {READER} from './action';
import ClsReader  from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';
import { scanner } from '../scanner/saga';

// eslint-disable-next-line prefer-const
let reader = new ClsReader('clsReader');
// reader.dispatcher = defaultReaderDispatcher;

function* readerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (reader.doValidate(data)){
      const respAction = reader.doDispatch(data);
      if (respAction) {
        yield put(respAction);
      }    } else {
      // do nothing
    }
  } catch (e) {
    CommonLog.Error(e)
  }
}

export default function* watchReader(): Saga<void> {
  yield takeLatest(READER.READ_NEW_DATA, readerHandler);
}
