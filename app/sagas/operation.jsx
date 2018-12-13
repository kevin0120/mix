import { select, put, take, call } from 'redux-saga/effects';
import {
  fetchRoutingWorkcenter,
  fetchWorkorder,
  jobManual,
  pset,
  fetchNextWorkOrder
} from './api/operation';
import { OPERATION, RUSH } from '../actions/actionTypes';
import { openShutdown } from '../actions/shutDownDiag';
import { OPERATION_RESULT, OPERATION_STATUS } from '../reducers/operations';
import { addNewStory, clearStories, STORY_TYPE } from './timeline';
import { toolEnable, toolDisable } from '../actions/tools';
import { setResultDiagShow } from '../actions/resultDiag';
import {
  fetchOngoingOperationOK,
  cleanOngoingOperation
} from '../actions/ongoingOperation';

// const lodash = require('lodash');

// 监听作业
export function* watchOperation() {
  while (true) {
    const action = yield take([
      OPERATION.STARTED,
      OPERATION.VERIFIED,
      OPERATION.FINISHED,
      OPERATION.RESET
    ]);
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
        if(state.setting.operationSettings.opMode === 'order'){
          yield call(getNextWorkOrderandShow);
        }
        yield put(setResultDiagShow(true));
        break;
      case OPERATION.RESET:
        yield put({ type: OPERATION.FINISHED });
        break;
      default:
        break;
    }
  }
}

function* getNextWorkOrderandShow() {
  try {
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
  } catch (e) {
    console.log(e);
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

    case OPERATION_STATUS.TIMEOUT:
      yield put({ type: OPERATION.PREDOING });
      break;

    default:
      break;
  }

  if (carID) {
    yield call(addNewStory, STORY_TYPE.INFO, source, carID);
    yield put({ type: OPERATION.TRIGGER.NEW_DATA, source, carID, carType: null });
  }

  if (carType) {
    yield call(addNewStory, STORY_TYPE.INFO, source, carType);
    yield put({ type: OPERATION.TRIGGER.NEW_DATA, source, carID: null, carType });
  }

  state = yield select();

  const triggers = state.setting.operationSettings.flowTriggers;

  let triggerFlagNum = 0;
  for (let i = 0; i < triggers.length; i += 1) {
    if (state.operations[triggers[i]] !== '') {
      triggerFlagNum += 1;
    }
  }

  if (triggerFlagNum === triggers.length || state.workMode.workMode === 'manual') {
    yield call(getOperation, job);
  }
}

// 定位作业
export function* getOperation(job) {
  try {
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
        try {
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
        } catch (e) {

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
      // yield put({ type: OPERATION.RESET });
    }
  } catch (e) {
    console.log(e);
  }

}

// 开始作业
export function* startOperation(data) {
  try {
    yield put(setResultDiagShow(false));

    let state = yield select();

    yield put({
      type: OPERATION.OPERATION.FETCH_OK,
      mode: state.setting.operationSettings.opMode,
      data
    });

    state = yield select();

    const { controllerMode } = state.workMode;

    const rushUrl = state.connections.masterpc;

    if (controllerMode === 'job') {
      // job模式

      const controllerSN = state.connections.controllers[0].serial_no;
      const {
        operationID,
        carType,
        carID,
        jobID,
        source
      } = state.operations;

      const { hmiSn } = state.setting.page.odooConnection;
      const userID = 1;
      const skip = false;
      let hasSet = false;
      if (state.setting.operationSettings.workMode === 'manual') {
        hasSet = true;
      }

      try {
        const resp = yield call(
          jobManual,
          rushUrl,
          controllerSN,
          carType,
          carID,
          userID,
          jobID,
          hmiSn.value,
          operationID,
          skip,
          hasSet,
          source
        );

        if (resp.status === 200) {
          // 程序号设置成功

          // 启动用具
          yield put(toolEnable());
          yield put({ type: OPERATION.STARTED });
        }
      } catch (e) {
        // 程序号设置失败
        yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
        // yield put({ type: OPERATION.RESET });
      }
    } else {
      yield put({ type: OPERATION.STARTED });
      yield call(doingOperation);
    }
  } catch (e) {

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

  const { operations } = state;

  let hasFail = false;
  let storyType = STORY_TYPE.PASS;

  for (let i = 0; i < data.length; i += 1) {
    if (data[i].result === OPERATION_RESULT.NOK) {
      hasFail = true;
      storyType = STORY_TYPE.FAIL;
    } else {
      storyType = STORY_TYPE.PASS;
    }

    const eti  = data[i].ti? data[i].ti.toString() : 'nil';

    const batch = `${(operations.activeResultIndex + 1).toString()}/${operations.results.length.toString()}`;

    yield call(
      addNewStory,
      storyType,
      `结果 ${batch}`,
      `T:${data[i].mi.toString()} W:${data[i].wi.toString()}`
    );
  }

  let rType = '';
  let continueDoing = false;

  if (hasFail) {
    if (
      operations.failCount + 1 >=
      operations.results[operations.activeResultIndex]
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
    operations.activeResultIndex + data.length >=
    operations.results.length
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
