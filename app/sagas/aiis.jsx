import { call, take, select } from 'redux-saga/effects';
import { ANDON } from '../actions/actionTypes';
import { OPERATION_STATUS } from '../reducers/operations';

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
    } else {
      // 空车信息
    }
  }
}
