// @flow
import {
  takeLatest,
  put,
  delay
} from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsRFID from './model';
import { RFID } from './action';
import { CommonLog } from '../../../../common/utils';
import type { tCommonActionType, tDeviceNewData } from '../type';

// eslint-disable-next-line prefer-const
let Rfid = new ClsRFID('RFID','demo rfid sn');

const DebounceWaitTime = 2000;

function* RFIDHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
  try {
    const { data } = action;
    if (Rfid.validate((data: string))) {
      const respAction = Rfid.doDispatch(data);
      if (respAction) {
        yield put(respAction);
      }
    }
  } catch (err) {
    CommonLog.lError(err);
  }
}

function* watchRFIDEvent(): Saga<void> {
  try {
    yield takeLatest(RFID.READ_NEW_DATA, RFIDHandler);
    yield delay(DebounceWaitTime);
  } catch (err) {
    CommonLog.lError(err);
  }
}

export default watchRFIDEvent;
