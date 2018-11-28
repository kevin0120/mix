import { select, put, take, call } from 'redux-saga/effects';
import {
  fetchRoutingWorkcenter,
  fetchWorkorder,
  jobManual,
  pset,
  fetchNextWorkOrder,
} from './api/operation';
import { OPERATION, RUSH } from '../actions/actionTypes';
import { openShutdown } from '../actions/shutDownDiag';
import { OPERATION_RESULT, OPERATION_STATUS } from '../reducers/operations';
import { addNewStory, clearStories, STORY_TYPE } from './timeline';
import { toolEnable, toolDisable } from '../actions/tools';
import { setResultDiagShow } from '../actions/resultDiag';
import {fetchOngoingOperationOK, cleanOngoingOperation} from '../actions/ongoingOperation';

const lodash = require('lodash');


// 监听作业
export function* watchOperation() {
  while (true) {
    const action = yield take([OPERATION.STARTED,OPERATION.VERIFIED, OPERATION.FINISHED]);
    const state = yield select();
    const { workMode } = state;

    switch (action.type) {
      case OPERATION.VERIFIED:
        yield call(startOperation, action.data);
        break;
      case OPERATION.STARTED:
        yield put(setResultDiagShow(false));
        yield put(cleanOngoingOperation());
        break;

      case OPERATION.FINISHED:
        if (workMode.controllerMode === 'job') {
          // 工具禁用
          yield put(toolDisable());
        }

        yield call(getNextWorkOrderandShow);
        yield put(setResultDiagShow(true));
        break;

      default:
        break;
    }
  }
}

function* getNextWorkOrderandShow() {
  const state = yield select();

  const rushUrl = state.connections.masterpc;
  const { workcenterCode } = state.connections;
  const resp = yield call(
    fetchNextWorkOrder,
    rushUrl,
    workcenterCode);
  if (resp.status === 200) {
    yield put(fetchOngoingOperationOK(resp.data));
  }
}

// 触发作业
export function* triggerOperation(carID, carType, job, source) {
  let state = yield select();
  switch (state.operations.operationStatus) {
    case OPERATION_STATUS.DOING:
      return;
    case OPERATION_STATUS.READY:
      yield call(clearStories);
      yield put({ type: OPERATION.PREDOING });
      break;

    default:
      break;
  }

  if (carID) {
    yield call(addNewStory, STORY_TYPE.INFO, source, carID);
    yield put({ type: OPERATION.TRIGGER.NEW_DATA, carID });
  }

  if (carType) {
    yield call(addNewStory, STORY_TYPE.INFO, source, carType);
    yield put({ type: OPERATION.TRIGGER.NEW_DATA, carType });
  }

  state = yield select();

  const triggers = state.setting.operationSettings.flowTriggers;

  let triggerFlagNum = 0;
  for (let i = 0; i < triggers.length; i += 1) {
    if (state.operations[triggers[i]] !== '') {
      triggerFlagNum += 1;
    }
  }

  if (triggerFlagNum === triggers.length) {
    yield call(getOperation, job);
  }
}

// 定位作业
export function* getOperation(job) {
  const state = yield select();

  const rushUrl = state.connections.masterpc;
  const { workcenterCode } = state.connections;

  let fetchOK = false;
  let resp = null;

  if (state.workMode.workMode === 'manual') {
    // 手动模式

    if (job) {
      resp = yield call(
        fetchRoutingWorkcenter,
        rushUrl,
        workcenterCode,
        null,
        job
      );
      if (resp.status === 200) {
        fetchOK = true;
      }
    }
  } else {
    if (state.setting.operationSettings.opMode === 'op') {
      // 作业模式

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
        resp = yield call(fetchWorkorder, rushUrl, workcenterCode, code);
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

  const { controllerMode, workMode } = state.workMode;

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
      results
    } = state.operations;

    const { hmiSn } = state.setting.page.odooConnection;
    const userID = 1;
    const skip = false;
    let hasSet = false;
    if (state.setting.operationSettings.workMode === 'manual') {
      hasSet = true;
    }

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

      // 启动用具
      yield put(toolEnable());
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
  const { controllerMode } = state.workMode;

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
      results[activeResultIndex].pset
    );

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

  if (operations.activeResultIndex >= operations.results.length - 1) {
    yield put({ type: OPERATION.FINISHED });
  } else {
    yield put({ type: OPERATION.CONTINUE });
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

  let hasFail = false;
  let storyType = STORY_TYPE.PASS;

  for (let i = 0; i < data.length; i += 1) {
    if (data[i].result === OPERATION_RESULT.NOK) {
      hasFail = true;
      storyType = STORY_TYPE.FAIL;
    } else {
      storyType = STORY_TYPE.PASS;
    }

    yield call(
      addNewStory,
      storyType,
      '结果',
      `扭矩:${data[i].mi.toString()} 角度:${data[i].wi.toString()}`
    );
  }

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

  yield put({ type: rType, data });
  if (continueDoing) {
    yield call(doingOperation);
  }
}
