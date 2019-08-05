// @flow

import { put, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import {CONTROLLER} from './action';
import ClsController  from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';

// eslint-disable-next-line prefer-const
const controller = new ClsController('controller', 'Dummy Serial Number');
const controller2 = new ClsController('controller', 'Dummy Serial Number2');

// reader.dispatcher = defaultReaderDispatcher;

function* controllerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (controller.doValidate(data)){
      yield put(controller.doDispatch(data))
    } else {
      // do nothing
    }
  } catch (e) {
    CommonLog.Error(e)
  }
}

export default function* watchController(): Saga<void> {
  yield takeLatest(CONTROLLER.CONTROLLER_SOCKET_DATA, controllerHandler);
}
