import { call, take, select, put } from 'redux-saga/effects';
import { ANDON, OPERATION } from '../actions/actionTypes';
import { OPERATION_SOURCE, OPERATION_STATUS } from '../reducers/operations';
import { triggerOperation } from './operation';
import { jobManual } from './api/operation';

export function* watchAiis() {
  while (true) {
    const { json } = yield take(ANDON.NEW_DATA);
    yield call(handleAiisData, json);
  }
}

export function* handleAiisData(data) {
  const state = yield select();
  if (state.operations.operationStatus !== OPERATION_STATUS.DOING) {
    if (data.vin_code.length) {
      // 车辆拧紧作业
      yield call(triggerOperation, data.vin_code, data.cartype_code, null, OPERATION_SOURCE.ANDON);
    } else {
      // 空车信息

      const {
        carType,
        carID,
        productID,
        workcenterID,
        results
      } = state.operations;

      const { jobID } = state.setting.operationSettings.emptyCarJob;

      const controllerSN = state.connections.controllers[0].serial_no;
      const rushUrl = state.connections.masterpc;
      const { hmiSn } = state.setting.page.odooConnection;
      const userID = 1;
      const skip = true;
      const hasSet = false;

      const resp = yield call(
        jobManual,
        rushUrl,
        controllerSN,
        carType,
        carID,
        userID,
        jobID,
        results,
        hmiSn.value,
        productID,
        workcenterID,
        skip,
        hasSet,
        ''
      );

      if (resp.statusCode !== 200) {
        yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
      }
    }
  }
}
