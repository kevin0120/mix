import { select, put, take, call } from 'redux-saga/effects';
import {
  fetchRoutingWorkcenter,
  fetchWorkorder,
  jobManual,
  pset
} from './api/operation';
import { OPERATION, RUSH, SHUTDOWN_DIAG } from '../actions/actionTypes';
import { openShutdown } from '../actions/shutDownDiag';
import { OPERATION_RESULT } from '../reducers/operations';
import { closeDiag, confirmDiag, openDiag } from './shutDownDiag';

const lodash = require('lodash');

// 监听作业
export function* watchOperation() {
  while (true) {
    const action = yield take([
      OPERATION.VERIFIED,
    ]);
    switch (action.type) {
      case OPERATION.VERIFIED:
        yield call(startOperation, action.data);
        break;
      default:
        break;
    }
  }
}

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
      const { preCheck } = state.setting.operationSettings;
      resp = e.response;
      if (resp.status === 409) {
        fetchOK = true;

        if (preCheck) {
          yield put(openShutdown('verify', resp.data));
          return;
        }
      }
    }
  }

  if (fetchOK) {
    // 定位作业成功
    // 开始作业

    yield call(startOperation, resp.data);
  } else {
    // 定位作业失败
    yield put({ type: OPERATION.OPERATION.FETCH_FAIL });
  }
}

// 开始作业
export function* startOperation(data) {

  let state = yield select();

  yield put({
    type: OPERATION.OPERATION.FETCH_OK,
    mode: state.setting.operationSettings.opMode,
    data
  });

  state = yield select();

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
    } = state.operations;

    const { hmiSn } = state.setting.page.odooConnection;
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

    const { masterpc } = state.connections;
    const { activeResultIndex, failCount, results } = state.operations;
    const userID = 1;

    const resp = yield call(
      pset,
      masterpc,
      results[activeResultIndex].controller_sn,
      results[activeResultIndex].gun_sn,
      results[activeResultIndex].id,
      failCount + 1,
      userID,
      results[activeResultIndex].pset);

    if (resp.status !== 200) {
      // 程序号设置失败
      yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
      return false;
    }
  }

  return true;
}

// 复位继续作业
export function* continueOperation() {
  const state = yield select();
  const { operations } = state;

  if (operations.activeResultIndex >= (operations.results.length - 1)) {
    yield put({type: OPERATION.FINISHED});
  } else {
    yield put({type: OPERATION.CONTINUE});
    yield call(doingOperation);
  }
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
  }

  yield put({type: rType, data});
  if (continueDoing) {
    yield call(doingOperation);
  }
}
