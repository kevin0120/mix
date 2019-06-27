/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import {  takeLatest, put, select } from 'redux-saga/effects';

import { CONNECTION } from '../connections/action';
import { RUSH } from '../rush/action';
import { RFID } from '../rfid/action';
import { SYSTEM_INIT } from './action';
import { setLedStatusReady, setModBusIO } from '../io/saga';

import { setNewNotification } from '../notification/action';
import { initIO } from '../io/action';
import { initAiis } from '../aiis/action';

// const lodash = require('lodash');

function* sysInit(action) {
  try {
    // const {section='all'} = action;
    const state = yield select();
    // const {connections} = state.setting.system;
    // const url = connections.masterpc.connection;
    // const controllers = connections.controllers.map(i => i.serial_no);

    // yield put(startHealthzCheck(url, controllers)); // 启动healthzcheck 定时器

    // 初始化rush
    yield put({ type: RUSH.INIT });

    // 初始化io
    const modbusConfig = state.setting.page.modbus;
    setModBusIO(modbusConfig);

    if (state.setting.systemSettings.modbusEnable) {
      yield put(initIO());
    }

    // 初始化rfid
    if (state.setting.systemSettings.rfidEnabled) {
      yield put({ type: RFID.INIT });
    }

    // 初始化aiis(andon)
    if (state.setting.systemSettings.andonEnable) {
      yield put(initAiis());
    }

    setLedStatusReady();
  } catch (e) {
    yield put(setNewNotification('error', e.toString()));
  }
}

// export function* fetchConnectionFlow(baseUrl, hmiSN, aiis) {
//   try {
//     const fullUrl = `${baseUrl}/hmi.connections/${hmiSN}`;
//
//     const resp = yield call(fetchConnectionInfo, fullUrl);
//
//     if (resp.status === 200) {
//       // yield put({ type: CONNECTION.FETCH_OK, data: resp.data });
//       const { masterpc, rfid, io, controllers, info } = resp.data;
//       const data =  {
//         masterpc: masterpc.connection ? masterpc.connection : '',
//         rfid: rfid.connection ? rfid.connection : '',
//         io: io.connection ? io.connection : '',
//         workcenterCode: info.workcenter_code ? info.workcenter_code : '',
//         rework_workcenter: info.qc_workcenter ? info.qc_workcenter : '',
//         aiis,
//         controllers: lodash.isArray(controllers) ? controllers : []
//       };
//       yield put({ type: USER_CONFIGS.SAVE, section:'connections', newConfigs: data }); // 修改配置中的信息
//     }
//
//   } catch (e) {
//     yield put(setNewNotification('error', e.toString()));
//   }
// }

export default function* sysInitFlow() {
  try {
    yield takeLatest(SYSTEM_INIT, sysInit);
    // yield call(fetchConnectionFlow);
  } catch (e) {
    yield put(setNewNotification('error', e.toString()));
  }
}
