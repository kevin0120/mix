/* eslint-disable no-empty-function */
import { put, fork, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import type { tRushWebSocketData, tWebSocketEvent, tRushData } from './type';
import { WEBSOCKET_EVENTS as wse } from './constants';
import notifierActions from '../Notifier/action';
import { CommonLog } from '../../common/utils';
import {
  toolNewResults,
  toolStatusChange
} from '../external/device/tools/saga';
import rushActions from './action';
import readerNewData from '../external/device/reader/saga';
import scannerNewData from '../external/device/scanner/saga';
import ioWSDataHandlers from '../external/device/io/handleWSData';
import { deviceStatus } from '../external/device';
import orderWSDataHandlers from '../order/handleWSData';

export default function* (payload) {
  try {
    const d = /(^[^"]*);(.*)/.exec(payload);
    const dataArray = d[1].split(';');
    const event: tWebSocketEvent = dataArray[0].split(':').slice(-1)[0];
    const json: tRushWebSocketData = JSON.parse(d[2]);
    const { sn, type, data, ...otherInfo } = json;
    CommonLog.Info(`rush message (${event})(${json.type})`, {
      event,
      type,
      sn,
      data: JSON.stringify(data),
      otherInfo: JSON.stringify(otherInfo)
    });
    yield put(rushActions.data(json));
    if (rushDataHandlers[event]) {
      yield fork(rushDataHandlers[event], json); // 异步处理rush数据
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush handleData' });
  }
}

const handleData = (dataHandlers) => function* _handleData(rushData: tRushData<any, any>): Saga<void> {
  try {
    const { type, data } = rushData;
    yield call(dataHandlers[type], data);
  } catch (e) {
    CommonLog.lError(e);
  }
};

const rushDataHandlers: {
  [tWebSocketEvent]: (tRushData<any, any>)=>void | Saga<void>
} = {
  * [wse.maintenance](data: tRushData<string, { name: string }>) {
    try {
      yield put(
        notifierActions.enqueueSnackbar(
          'Maintenance',
          `新维护请求: ${data.type},${data.data.name}`
        )
      );
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event maintenance' });
    }
  },
  [wse.io]: handleData(ioWSDataHandlers),
  [wse.scanner]: scannerNewData,
  [wse.reader]: readerNewData,
  [wse.result]: toolNewResults,
  [wse.tool]: toolStatusChange,
  [wse.device]: deviceStatus,
  [wse.order]: handleData(orderWSDataHandlers)
};
