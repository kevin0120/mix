// @flow

import { call, takeLatest } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { READER } from './action';
import ClsReader from './model';
import type { tCommonActionType, tDeviceNewData } from '../../common/type';
import { CommonLog } from '../../common/utils';
import { symReader, AppendNewDevices } from '../global';
import type { tRushWebSocketData } from '../rush/type';
import type { tReaderData } from './model';


export const reader = new ClsReader(symReader);
AppendNewDevices(symReader, reader);
// TODO: 收到读卡器uid后，分发给user模块进行用户认证
// reader.dispatcher = defaultReaderDispatcher;


export default function* readerNewData(data: tRushWebSocketData): Saga<void> {
  try {
    const d = (data.data: tReaderData);
    CommonLog.Info(` Reader receive data: ${d.uid}`);
    // yield put(ReaderNewData(d.uid));
    yield call(reader.doDispatch, data);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event reader' });
  }
}
