/* eslint-disable no-lonely-if,camelcase */
import { select, put, take, call, fork, takeLatest } from 'redux-saga/effects';
import {
  fetchRoutingWorkcenter,
  fetchWorkorder,
  jobManual,
  pset,
  fetchNextWorkOrder,
  ak2Api
} from '../../api/operation';

import {
  OPERATION_RESULT,
  OPERATION_SOURCE,
  OPERATION_STATUS
} from './model';
import { addNewStory, clearStories } from '../timeline/saga';
import { STORY_TYPE } from '../timeline/model';
import { toolEnable, toolDisable, TOOLS } from '../tools/action';
import { setResultDiagShow } from '../resultDiag/action';
import {
  switch2Ready,
  switch2PreDoing,
  operationConflictDetected,
  operationTriggerBlock,
  OPERATION
} from './action';
import {
  fetchOngoingOperationOK,
  cleanOngoingOperation
} from '../ongoingOperation/action';
import { lError, Info } from '../../logger';
import { setNewNotification } from '../notification/action';
import { watchWorkers } from '../util';
import configs from '../../shared/config';

const lodash = require('lodash');

const operationMainProgress = {
  [OPERATION.STARTED]: [call, operationStarted],
  [OPERATION.FINISHED]: [call, operationFinished],
  [OPERATION.RESET]: [put, switch2Ready()]
};

const operationWorkers = {
  // bypass
  [OPERATION.BYPASS.CONFIRM]: [fork, bypassConfirm],
  // conflict
  [OPERATION.CONFLICT.DETECTED]: [fork, conflictDetected],
  [OPERATION.CONFLICT.CONFIRM]: [fork, startOperation],
  [OPERATION.CONFLICT.CANCEL]: [fork, conflictCanceled]
};

export function* operationFlow() {
  try {
    yield fork(watchWorkers(operationMainProgress));
    yield fork(watchWorkers(operationWorkers));
    yield takeLatest(OPERATION.TRIGGER.TRIGGER, triggerOperation);
  } catch (e) {
    console.error(e);
  }
}

function* operationStarted() {
  try {
    yield put(setResultDiagShow(false));
    yield put(cleanOngoingOperation());
  } catch (e) {
    console.error(e);
  }
}

function* operationFinished(action) {
  try {
    const { showDiag } = action;
    yield put(operationTriggerBlock(false));

    const state = yield select();


    if (state.setting.operationSettings.opMode === 'order') {
      yield call(getNextWorkOrderandShow);
    }
    if (showDiag) {
      yield put(setResultDiagShow(true));
      // 工具禁用

      yield put(toolDisable('作业结束'));

    }
  } catch (e) {
    console.error(e);
  }
}

function* getNextWorkOrderandShow() {
  try {
    const state = yield select();

    const { rush: rushUrl, workcenterCode } = state.connections;
    const resp = yield call(fetchNextWorkOrder, rushUrl, workcenterCode);
    if (resp.status === 200) {
      yield put(fetchOngoingOperationOK(resp.data));
    }
  } catch (e) {
    console.error(`getNextWorkOrderandShow ${e.message}`);
  }
}

// 触发作业
function* triggerOperation(action) {
  try {

    const { carID, carType, job, source } = action;
    const rState = yield select();

    if (rState.router.location.pathname !== '/working') {
      // 不在作业页面，直接返回
      return;
    }

    if (rState.operations.trigger.block && source === 'RFID') {
      return;
    }

    switch (rState.operations.operationStatus) {
      case OPERATION_STATUS.DOING:
        return;
      case OPERATION_STATUS.READY:
        yield call(clearStories);
        yield put(switch2PreDoing());
        break;
      case OPERATION_STATUS.TIMEOUT:
        yield put(switch2PreDoing());
        break;
      default:
        break;
    }

    yield put({ type: OPERATION.SOURCE.SET, source });

    if (carID) {
      yield call(addNewStory, STORY_TYPE.INFO, source, carID);
      yield put({
        type: OPERATION.TRIGGER.NEW_DATA,
        carID,
        carType: null
      });
    }

    if (carType) {
      yield call(addNewStory, STORY_TYPE.INFO, source, carType);
      yield put({
        type: OPERATION.TRIGGER.NEW_DATA,
        carID: null,
        carType
      });
    }

    if (
      source === OPERATION_SOURCE.SCANNER &&
      rState.workMode.workMode === 'manual'
    ) {
      // 手动模式下,扫码枪接收到讯息,不获取作业以及切换工作状态
      return;
    }


    const triggers = rState.workMode.workMode === 'manual' ? ['carID'] : rState.setting.operationSettings.flowTriggers; // 手动模式下，只需要车辆信息即可触发作业

    const operations = yield select(state => state.operations);

    let triggerFlagNum = 0;
    for (let i = 0; i < triggers.length; i += 1) {
      if (operations[triggers[i]] !== '') {
        triggerFlagNum += 1;
      }
    }

    if (triggerFlagNum === triggers.length) {
      yield call(getOperation, job);
    }
  } catch (e) {
    console.error(e);
  }
}

// 定位作业
export function* getOperation(job) {
  try {
    const state = yield select();

    const { rush: rushUrl, workcenterCode } = state.connections;

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

        // const hmiSN = state.setting.page.odooConnection.hmiSn.value;
        const code = state.operations.carID;
        try {
          resp = yield call(fetchWorkorder, rushUrl, workcenterCode, code);
          if (resp.status === 200) {
            fetchOK = true;
          }
        } catch (e) {
          const { preCheck = false } = state.setting.operationSettings;
          resp = e.response;
          if (resp.status === 409) {
            fetchOK = true;

            if (preCheck) {
              yield put(operationConflictDetected(resp.data));
              return;
            }
          } else {
            yield put(setNewNotification('Error', `获取工单失败:${e.message}`, {
              workMode: state.workMode.workMode,
              opMode: state.setting.operationSettings.opMode,
              carID: code,
              response: resp
            }));
            yield put({ type: OPERATION.OPERATION.FETCH_FAIL });
            yield call(clearStories);
          }
        }
      }
    }

    if (fetchOK) {
      // 定位作业成功
      // 开始作业
      yield call(startOperation, { data: resp.data });
    } else {
      // 定位作业失败
      yield put(setNewNotification('Error', '定位作业失败', {
        workMode: state.workMode.workMode,
        opMode: state.setting.operationSettings.opMode,
        carID: state.operations.carID,
        response: resp
      }));
      yield put({ type: OPERATION.OPERATION.FETCH_FAIL });
      yield call(clearStories);
      // yield put({ type: OPERATION.RESET });
    }
  } catch (e) {
    const state = yield select();
    yield put(setNewNotification('Error', `获取作业失败:${e.message}`, {
      workMode: state.workMode.workMode,
      opMode: state.setting.operationSettings.opMode,
      carID: state.operations.carID
    }));
    yield put({ type: OPERATION.OPERATION.FETCH_FAIL });
    yield call(clearStories);
  }
}

// 开始作业
export function* startOperation(action) {
  try {
    yield put(setResultDiagShow(false));

    const { data } = action;
    yield put({
      type: OPERATION.OPERATION.FETCH_OK,
      mode: configs.operationSettings.opMode,
      data
    });

    const state = yield select();

    const { controllerMode } = state.workMode;

    const rushUrl = state.connections.rush;

    if (controllerMode === 'job') {
      // job模式

      // const controllerSN = state.connections.controllers[0].serial_no;
      const { operationID, carType, carID, jobID, source, results } = state.operations;

      const { hmiSn } = state.setting.page.odooConnection;

      // const toolSN = state.setting.systemSettings.defaultToolSN || "";

      const { controller_sn, gun_sn } = lodash.reduce(results, (result, value) => {
        if (result.controller_sn && value.controller_sn !== result.controller_sn) {
          console.error('结果中的controller_sn不匹配');
        }
        if (result.gun_sn && value.gun_sn !== result.gun_sn) {
          console.error('结果中的gun_sn不匹配');
        }
        return {
          controller_sn: value.controller_sn || result.controller_sn,
          gun_sn: value.gun_sn || result.gun_sn
        };
      }, {});

      const userID = 1;
      const skip = false;
      let hasSet = false;
      if (state.workMode.workMode === 'manual') {
        hasSet = true;
      }
      try {
        const resp = yield call(
          jobManual,
          rushUrl,
          controller_sn,
          gun_sn,
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

          // 设置workorder_id
          yield put({
            type: OPERATION.JOB_MANUAL.OK,
            workorderID: resp.data.workorder_id
          });
          // 启动用具
          yield put(toolEnable('开始作业'));
          yield put({ type: OPERATION.STARTED });
        }
      } catch (e) {
        // 程序号设置失败
        console.error(e);
        yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
        // yield put({ type: OPERATION.RESET });
      }
    } else {
      const rt = yield call(doingOperation, controllerMode);
      if (rt) {
        yield put({ type: OPERATION.STARTED });
      }
    }
  } catch (e) {
    console.error(`startOperation ${e.message}`);
  }
}

// 处理作业过程
export function* doingOperation(controllerMode) {
  if (controllerMode === 'pset') {
    // pset模式
    try {
      const state = yield select();
      const { rush } = state.connections;
      const {
        activeResultIndex,
        failCount,
        results,
        workorderID
      } = state.operations;
      const userID = 1;
      yield call(
        pset,
        rush,
        results[activeResultIndex].controller_sn,
        results[activeResultIndex].gun_sn,
        0,
        failCount + 1,
        userID,
        results[activeResultIndex].pset,
        workorderID,
        results[activeResultIndex].group_sequence
      );
    } catch (e) {
      // 程序号设置失败
      yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
      yield put(setNewNotification('Error', 'pset failed'), {
        workMode: state.workMode.workMode,
        opMode: state.setting.operationSettings.opMode,
        carID: state.operations.carID
      });
      return false;
    }
    return true;
  }
}

// 复位继续作业
export function* continueOperation() {
  try {
    const state = yield select();
    const { operations } = state;
    const { controllerMode } = state.workMode;

    // if (state.operations.operationStatus === OPERATION_STATUS.FAIL) {
    //
    // }

    if (operations.activeResultIndex >= operations.results.length - 1) {

      yield put(switch2Ready());
    } else {
      yield put({ type: OPERATION.CONTINUE });
      yield call(doingOperation, controllerMode);
    }
  } catch (e) {
    console.error(`continueOperation fail: ${e.message}`);
  }
}

// 监听结果
export function* watchResults() {
  try {
    while (true) {
      const { data } = yield take(RUSH.NEW_RESULTS);
      yield call(handleResults, data);
    }
  } catch (e) {
    console.error(`watchResults: ${e.message}`);
  }
}

// 处理结果
export function* handleResults(data) {
  try {
    const state = yield select();

    const { operations } = state;

    const { controllerMode } = state.workMode;

    let hasFail = false;
    let storyType = STORY_TYPE.PASS;

    const batch = `${(
      operations.activeResultIndex + 1
    ).toString()}/${operations.results[
    operations.results.length - 1
      ].group_sequence.toString()}`;

    for (let i = 0; i < data.length; i += 1) {
      if (data[i].result === OPERATION_RESULT.NOK) {
        hasFail = true;
        storyType = STORY_TYPE.FAIL;
      } else if (data[i].result === OPERATION_RESULT.OK) {
        storyType = STORY_TYPE.PASS;
      } else {
        yield call(
          addNewStory,
          STORY_TYPE.FAIL,
          `结果 ${batch}`,
          `执行策略 ${data[i].result}`
        );
        return;
      }

      // const eti  = data[i].ti? data[i].ti.toString() : 'nil';

      yield call(
        addNewStory,
        storyType,
        `结果 ${batch}`,
        `T=${data[i].mi.toString()}Nm A=${data[i].wi.toString()}°`
      );
    }

    let rType = '';
    let continueDoing = false;
    const payLoad = { data };
    if (hasFail) {
      if (
        operations.failCount + 1 >=
        operations.results[operations.activeResultIndex].max_redo_times
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
      payLoad.showDiag = true;
    } else {
      // 继续作业
      rType = OPERATION.RESULT.OK;
      continueDoing = true;
    }

    yield put({ type: rType, ...payLoad });
    if (continueDoing) {
      yield call(doingOperation, controllerMode);
    }
  } catch (e) {
    console.error(`handleResults: ${e.message}`);
  }
}

export function* ak2() {
  try {
    const { workorderID, results, activeResultIndex, failCount } = yield select(
      state => state.operations
    );
    const { rush } = yield select(state => state.connections);
    yield call(
      ak2Api,
      rush,
      results[activeResultIndex].controller_sn,
      results[activeResultIndex].gun_sn,
      workorderID,
      results[activeResultIndex].pset,
      results[activeResultIndex].group_sequence,
      failCount + 1
    );
  } catch (e) {
    console.error(`ak2: ${e.message}`);
  }
}

function* bypassConfirm() {
  try {
    const state = yield select();
    const { operations: op } = state;
    const { carID } = op;
    const { enableAk2 = true } = state.setting.operationSettings;
    if (enableAk2) {
      yield call(ak2);
    }
    Info(`车辆已放行 车辆ID:${carID}`);

    yield put(switch2Ready());
  } catch (e) {
    console.error(e);
  }
}

function* conflictDetected(action) {
  try {
    yield put(operationTriggerBlock(true));
    const state = yield select();
    const { data } = action;
    const { enableConflictOP = false } = state.setting.systemSettings;
    if (!enableConflictOP) {
      yield put(setNewNotification('Warn', `设定为不允许重复拧紧同一张工单 VIN: ${data.vin}`));
      // return; // 直接返回, 不关闭模式对话框
    }
  } catch (e) {
    console.error(e);
  }
}

function* conflictCanceled() {
  try {
    yield put(operationTriggerBlock(false));

    yield put(switch2Ready());
  } catch (e) {
    console.error(e);
  }
}
