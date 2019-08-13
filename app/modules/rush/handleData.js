/* eslint-disable no-empty-function */
import { put, fork } from 'redux-saga/effects';
import type { tRushWebSocketData, tWebSocketEvent } from './type';
import { WEBSOCKET_EVENTS as wse } from './type';
import { setNewNotification } from '../notification/action';
import { CommonLog } from '../../common/utils';
import { toolNewResults, toolStatusChange } from '../external/device/tools/saga';
import rushActions from './action';
import readerNewData from '../external/device/reader/saga';
import scannerNewData from '../external/device/scanner/saga';
import ioNewData from '../external/device/io/saga';
import { deviceStatus } from '../external/device';

export default function* (payload) {
  try {
    const dataArray = payload.split(';');
    const event: tWebSocketEvent = dataArray[0].split(':').slice(-1)[0];

    const json: tRushWebSocketData = JSON.parse(dataArray.slice(-1));
    if (rushDataHandlers[event]) {
      yield fork(rushDataHandlers[event], json); // 异步处理rush数据
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush handleData' });
  }
}

const rushDataHandlers = {
  * [wse.maintenance](data: tRushWebSocketData) {
    try {
      yield put(setNewNotification('Maintenance', `新维护请求: ${data.type},${data.data.name}`));

    } catch (e) {
      CommonLog.lError(e, { at: 'rush event maintenance' });
    }
  },

  * [wse.job]() {

  },

  * [wse.odoo]() {

  },


  [wse.controller](data: tRushWebSocketData) {
    try {
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event controller' });
    }
  },


  * [wse.tightening_device](data: tRushWebSocketData) {
    try {
      // 初始化所有拧紧设备

    } catch (e) {
      CommonLog.lError(e, { at: 'rush event tightening_device' });
    }
  },
  * [wse.reply](data: tRushWebSocketData) {
    try {
      yield put(rushActions.reply(data));
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event reply' });
    }
  },
  [wse.io]: ioNewData,

  [wse.scanner]: scannerNewData,

  [wse.reader]: readerNewData,

  [wse.result]: toolNewResults,

  [wse.tool]: toolStatusChange,

  [wse.device]: deviceStatus
};
