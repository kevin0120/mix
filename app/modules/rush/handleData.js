// @flow
/* eslint-disable no-empty-function */
import { put, fork, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import type { tWebSocketEvent, tRushData } from './type';
import { WEBSOCKET_EVENTS as wse, WS_RUSH } from './constants';
import notifierActions from '../Notifier/action';
import { CommonLog } from '../../common/utils';
import { toolNewResults, toolStatusChange } from '../device/tools/saga';
import rushActions from './action';
import readerNewData from '../device/reader/saga';
import scannerNewData from '../device/scanner/saga';
import ioWSDataHandlers from '../device/io/handleWSData';
import { deviceStatus } from '../deviceManager/handlerWSData';
import orderWSDataHandlers from '../order/handleWSData';
import systemInfoActions from '../systemInfo/action';
import healthzActions from '../healthz/action';

const registerHandlers = {
  * [WS_RUSH.RUSH_DATA]({ workcenter }) {
    try {
      // const { workcenter } = data;
      if (workcenter) {
        yield put(systemInfoActions.setWorkcenter(workcenter));
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};

export default function* (payload: string): Saga<void> {
  try {
    if (!payload) {
      return;
    }
    const d = /(^[^"]*);(.*)/.exec(payload);
    if (!d) {
      return;
    }
    const dataArray = d[1].split(';');
    const event: tWebSocketEvent = dataArray[0].split(':').slice(-1)[0];
    const json: tRushData<string, any> = JSON.parse(d[2]);
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

const handleData = dataHandlers =>
  function* _handleData(rushData: tRushData<any, any>): Saga<void> {
    try {
      const { type, data } = rushData;
      const handler = dataHandlers[type];
      if (handler) {
        yield call(handler, data);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  };

const rushDataHandlers: {
  [tWebSocketEvent]: (tRushData<any, any>) => void | Saga<void>
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
  [wse.order]: handleData(orderWSDataHandlers),
  [wse.register]: handleData(registerHandlers),
  * [wse.odoo]({ type, data }) {
    try {
      if (type === 'WS_ODOO_STATUS') {
        const { name, status } = data;
        yield call(updateExternalStatus, name, status);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [wse.aiis]({ type, data }) {
    try {
      if (type === 'WS_AIIS_STATUS') {
        const { name, status } = data;
        yield call(updateExternalStatus, name, status);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [wse.exsys]({ type, data }) {
    try {
      if (type === 'WS_EXSYS_STATUS') {
        const { name, status } = data;
        yield call(updateExternalStatus, name, status);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [wse.controller]() {
    // 控制器状态推送在设备列表中处理
  },
  * [wse.debug]({ type, data }) {
    try {
      yield put(notifierActions.enqueueSnackbar('test message'));
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};

function* updateExternalStatus(name, status) {
  try {
    const statusBool = status === 'online';
    yield put(healthzActions.data({ [name]: statusBool }));
  } catch (e) {
    CommonLog.lError(e, { at: 'updateExternalStatus' });
  }
}
