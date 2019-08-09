/* eslint-disable no-empty-function */
import { put, fork } from 'redux-saga/effects';
import type { tBarcode, tReader, tRushWebSocketData, tWebSocketEvent } from './type';
import { WEBSOCKET_EVENTS as wse } from './type';
import { setNewNotification } from '../notification/action';
import { CommonLog } from '../../common/utils';
import type { tIOContact, tIODirection, tIOWSMsgType } from '../io/type';
import { onchangeIO } from '../io/action';
import { toolNewResults, toolStatusChange } from '../tools/action';
import rushActions from './action';
import readerNewData from '../reader/saga';
import scannerNewData from '../scanner/saga';
import ioNewData from '../io/saga';

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


  * [wse.result](data: tRushWebSocketData) {
    try {
      CommonLog.Info(` tool new results: ${data.data}`);
      yield put(toolNewResults(data.data));
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event result' });
    }
  },


  [wse.controller](data: tRushWebSocketData) {
    try {
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event controller' });
    }
  },

  * [wse.tool](data: tRushWebSocketData) {
    try {
      yield put(toolStatusChange(data.tool_sn, data.status, data.reason));
    } catch (e) {
      CommonLog.lError(e, { at: 'rush event tool' });
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

  [wse.reader]: readerNewData
};
