// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsReader from './model';
import { CommonLog } from '../../../../common/utils';
import type { tReaderData } from './model';


// TODO: 收到读卡器uid后，分发给user模块进行用户认证

type tReaderRushData = {
  type: string,
  data: tReaderData
}

export default function* readerNewData(data: tReaderRushData): Saga<void> {
  try {
    const d = (data.data: tReaderData);
    CommonLog.Info(` Reader receive data: ${d.uid}`);
    // yield put(ReaderNewData(d.uid));
    yield call(reader.doDispatch, data);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event reader' });
  }
}
