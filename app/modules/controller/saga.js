// @flow

import { call, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CONTROLLER } from './action';
import ClsController from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';
import type { tRushWebSocketData } from '../rush/type';

// eslint-disable-next-line prefer-const
const controller = new ClsController('controller', 'Dummy Serial Number');

export default function* controllerNewData(data: tRushWebSocketData): Saga<void> {
  try {
    yield call(controller.doDispatch,data);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event controller' });
  }
}
// function* controllerHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
//   try {
//     const { data } = action;
//     if (controller.doValidate(data)){
//       yield put(controller.doDispatch(data))
//     } else {
//       // do nothing
//     }
//   } catch (e) {
//     CommonLog.lError(e)
//   }
// }
//
// export default function* watchController(): Saga<void> {
//   yield takeLatest(CONTROLLER.CONTROLLER_SOCKET_DATA, controllerHandler);
// }
