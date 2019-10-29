// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../../../common/utils';
import { getDevicesByType } from '../../../deviceManager/devices';
import type { tReaderData } from './typeDef';
import { deviceType } from '../../../deviceManager/constants';

type tReaderRushData = {
  type: string,
  data: tReaderData
};

export default function* readerNewData(data: tReaderRushData): Saga<void> {
  try {
    const d = (data.data: tReaderData);
    CommonLog.Info(` Reader receive data: ${d.uid}`);
    // yield put(ReaderNewData(d.uid));
    const readers = getDevicesByType(deviceType.reader);
    if (readers && readers.length > 0) {
      yield call(readers[0].doDispatch, data);
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event reader' });
  }
}
