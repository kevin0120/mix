import { select, put, take, call } from 'redux-saga/effects';
import {
  fetchRoutingWorkcenter,
  fetchWorkorder,
  jobManual
} from './api/operation';
import { OPERATION, RUSH } from '../actions/actionTypes';
import { openShutdown } from '../actions/shutDownDiag';
import { OPERATION_RESULT } from '../reducers/operations';

const lodash = require('lodash');


// 定位作业
export function* getOperation() {
  const state = yield select();

  const rushUrl = state.connections.masterpc;

  let fetchOK = false;
  let resp = null;
  if (state.setting.operationSettings.opMode === 'op') {
    // 作业模式

    const { workcenterCode } = state.connections;
    const { carType } = state.operations;
    resp = yield call(
      fetchRoutingWorkcenter,
      rushUrl,
      workcenterCode,
      carType,
      null
    );
    if (resp.status === 200) {
      fetchOK = true;
    }
  } else {
    // 工单模式
    const hmiSN = state.setting.page.odooConnection.hmiSn.value;
    const code = state.operations.carID;
    try {
      resp = yield call(fetchWorkorder, rushUrl, hmiSN, code);
      if (resp.status === 200) {
        fetchOK = true;
      }
    } catch (e) {
      if (e.response.status === 409) {
        yield put(openShutdown, 'verify', e.response.data);
      }
    }
  }

  if (fetchOK) {
    // 定位作业成功
    yield put({
      type: OPERATION.OPERATION.FETCH_OK,
      mode: state.setting.operationSettings.opMode,
      data: resp.data
    });

    // 开始作业
    yield call(startOperation);
  } else {
    // 定位作业失败
    yield put({ type: OPERATION.OPERATION.FETCH_FAIL });
  }
}

// 开始作业
export function* startOperation() {
  const state = yield select();
  const { controllerMode, workMode } = state.setting.operationSettings;

  const rushUrl = state.connections.masterpc;

  if (controllerMode === 'job') {
    // job模式

    const controllerSN = state.connections.controllers[0].serial_no;
    const {
      carType,
      carID,
      productID,
      workcenterID,
      jobID,
      results,
      hmiSn
    } = state.operations;
    const userID = 1;
    const skip = false;
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
      workMode
    );

    if (resp.status === 200) {
      // 程序号设置成功
      yield put({ type: OPERATION.STARTED });
    } else {
      // 程序号设置失败
      yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
    }
  } else {
    const rt = yield call(doingOperation);
    if (rt) {
      yield put({ type: OPERATION.STARTED });
    }
  }
}

// 处理作业过程
export function* doingOperation() {
  const state = yield select();
  const { controllerMode } = state.setting.operationSettings;

  if (controllerMode === 'pset') {
    // pset模式
  }

  return true;
}

// 监听结果
export function* watchResults() {
  while (true) {
    const { data } = yield take(RUSH.NEW_RESULTS);
    yield call(handleResults, data);
  }
}

// 处理结果
export function* handleResults(data) {
  const state = yield select();

  const hasFail = lodash.every(data, v => v.result === OPERATION_RESULT.NOK);


  let rType = '';
  let continueDoing = false;

  if (hasFail) {
    if (
      state.operations.failCount + 1 >=
      state.operations.results[state.operations.activeResultIndex]
        .max_redo_times
    ) {
      // 作业失败
      rType = OPERATION.FAILED;
    } else {
      // 重试
      rType = OPERATION.RESULT.NOK;
      continueDoing = true;
    }
  } else if (
    state.operations.activeResultIndex + data.length >=
    state.operations.results.length
  ) {
    // 作业完成
    rType = OPERATION.FINISHED;
  } else {
    // 继续作业
    rType = OPERATION.RESULT.OK;
    continueDoing = true;

  yield put({type: rType, data});
  if (continueDoing) {
    yield call(doingOperation);
  }
}
