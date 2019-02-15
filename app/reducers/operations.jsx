import { OPERATION } from '../actions/actionTypes';

import {
  setLedStatusDoing,
  setLedError,
  setLedStatusReady,
  sOn
} from '../sagas/io';

import sortObj from '../common/utils';
import { genReducers } from './utils';

export const OPERATION_STATUS = {
  READY: 'Ready',
  PREDOING: 'PreDoing',
  DOING: 'Doing',
  TIMEOUT: 'Timeout',
  FAIL: 'Fail'
};

export const OPERATION_RESULT = {
  OK: 'OK',
  NOK: 'NOK',
  LN: 'LSN'
};

export const OPERATION_SOURCE = {
  SCANNER: 'SCANNER',
  RFID: 'RFID',
  ANDON: 'ANDON',
  MANUAL: 'MANUAL'
};

const defaultOperations = {
  workorderID: 0,
  operationID: 0,
  operationStatus: 'Ready',
  carID: '',
  carType: '',
  activeResultIndex: 0,
  failCount: 0, // 失敗的次數
  jobID: 0,
  maxOpTimes: 0,
  workSheet: null,
  productID: -1,
  workcenterID: -1,
  lnr: '',
  source: '',
  bypassToConfirm: false,
  conflict: {
    toConfirm: false,
    data: {}
  },
  trigger: {
    block: false
  },
  results: [
    // {
    //   id: -1,
    //   controller_sn: '',
    //   gun_sn: '',
    //   pset: -1,
    //   max_redo_times: 3,
    //   offset_x: 0,
    //   offset_y: 0,
    //   sequence: 0,
    //   group_sequence: 0,
    //   ti: 0,
    //   mi: 0,
    //   wi: 0,
    //   result: '',
    //   batch:''
    // }
  ]
};

type actionType = {
  +type: string,
  +data: object,
  +carID: string,
  +carType: string,
  +source: string
};

const reducers = {
  [OPERATION.TRIGGER.NEW_DATA]: newTriggerData,
  [OPERATION.SOURCE.SET]: setSource,
  [OPERATION.OPERATION.FETCH_OK]: newOperation,
  [OPERATION.JOB_MANUAL.OK]: setWorkorderID,
  [OPERATION.OPERATION.FETCH_FAIL]: operationSwitchReady,
  [OPERATION.STARTED]: operationStarted,
  [OPERATION.RESULT.OK]: operationResultOK,
  [OPERATION.RESULT.NOK]: operationResultNOK,
  [OPERATION.FAILED]: operationFailed,
  [OPERATION.FINISHED]: operationFinished,
  [OPERATION.CONTINUE]: operationContinue,
  [OPERATION.PREDOING]: switchOperationPredoing,
  [OPERATION.TIMEOUT]: switchOperationTimeout,
  [OPERATION.BYPASS.IO]: operationBypassIO,
  [OPERATION.BYPASS.CONFIRM]: operationBypassConfirmed,
  [OPERATION.BYPASS.CANCEL]: operationBypassCancel,
  [OPERATION.CONFLICT.CONFIRM]: operationConflictConfirm,
  [OPERATION.CONFLICT.DETECTED]: operationConflictDetected,
  [OPERATION.CONFLICT.CANCEL]: operationConflictCanceled,
  [OPERATION.TRIGGER.BLOCK]: operationTriggerBlocker
};

export default genReducers(reducers, defaultOperations);

// export default function operations(
//   state: object = defaultOperations,
//   action: actionType
// ) {
//   switch (action.type) {
//     case OPERATION.TRIGGER.NEW_DATA:
//       return newTriggerData(state, action.carID, action.carType);
//     case OPERATION.SOURCE.SET:
//       return setSource(state, action.source);
//     case OPERATION.OPERATION.FETCH_OK:
//       return newOperation(state, action.mode, action.data);
//     case OPERATION.JOB_MANUAL.OK:
//       return setWorkorderID(state, action.workorderID);
//     case OPERATION.OPERATION.FETCH_FAIL:
//       return operationSwitchReady(state);
//     case OPERATION.STARTED:
//       return operationStarted(state);
//     case OPERATION.RESULT.OK:
//       return operationResultOK(state, action.data);
//     case OPERATION.RESULT.NOK:
//       return operationResultNOK(state, action.data);
//     case OPERATION.FAILED:
//       return operationFailed(state, action.data);
//     case OPERATION.FINISHED:
//       return operationFinished(state, action.data);
//     case OPERATION.CONTINUE:
//       return operationContinue(state);
//     case OPERATION.PREDOING:
//       return switchOperationPredoing(state, OPERATION_STATUS.PREDOING);
//     case OPERATION.TIMEOUT:
//       return switchOperationTimeout(state, OPERATION_STATUS.TIMEOUT);
//     case OPERATION.BYPASS.IO:
//       return operationBypassIO(state);
//     case OPERATION.BYPASS.CONFIRM:
//       return operationBypassConfirmed(state);
//     case OPERATION.CONFLICT.CONFIRM:
//       return operationConflictConfirm(state);
//     case OPERATION.CONFLICT.DETECTED:
//       return operationConflictDetected(state,action);
//       case OPERATION.CONFLICT.CANCEL:
//       return operationConflictCanceled(state,action);
//     default:
//       return state;
//   }
// }

function operationTriggerBlocker(state, action) {
  const { block } = action;
  return {
    ...state,
    trigger: {
      ...state.trigger,
      block
    }
  };
}

function operationConflictCanceled(state, action) {
  return {
    ...state,
    conflict: {
      data: action.data,
      toConfirm: false
    }
  };
}

function operationConflictDetected(state, action) {
  return {
    ...state,
    conflict: {
      data: action.data,
      toConfirm: true
    }
  };
}

function operationConflictConfirm(state) {
  return {
    ...state,
    conflict: {
      data: {},
      toConfirm: false
    }
  };
}

function operationBypassConfirmed(state) {
  return {
    ...state,
    bypassToConfirm: false
  };
}

function operationBypassCancel(state) {
  return {
    ...state,
    bypassToConfirm: false
  };
}

function operationBypassIO(state) {
  return {
    ...state,
    bypassToConfirm: true
  };
}

function newTriggerData(state, action) {
  const { carID, carType } = action;
  return {
    ...state,
    carID: carID !== null ? carID : state.carID,
    carType: carType !== null ? carType : state.carType
  };
}

function setSource(state, action) {
  const { source } = action;
  return {
    ...state,
    source
  };
}

function newOperation(state, action) {
  const { mode, data } = action;
  if (mode === 'op') {
    // 作业模式
    return {
      ...state,
      operationID: data.id,
      jobID: data.job,
      carType: data.product_type,
      maxOpTimes: data.max_op_time,
      workSheet: data.img,
      productID: data.product_id,
      workcenterID: data.workcenter_id,
      activeResultIndex: 0,
      results: data.points.sort((a, b) => a.group_sequence - b.group_sequence)
    };
  }

  // 工单模式
  return {
    ...state,
    jobID: data.job_id,
    carType: data.model,
    maxOpTimes: data.max_op_time,
    workSheet: data.work_sheet,
    results: data.results.sort((a, b) => a.group_sequence - b.group_sequence),
    activeResultIndex: 0,
    lnr: data.lnr,
    workorderID: data.workorder_id
  };
}

function setWorkorderID(state, action) {
  const { workorderID } = action;
  return {
    ...state,
    workorderID
  };
}

function operationSwitchReady(state) {
  setLedStatusReady();

  return {
    ...state,
    operationStatus: OPERATION_STATUS.READY,
    carID: '',
    carType: '',
    lnr: '',
    maxOpTimes: 0,
    failCount: 0,
    results: []
  };
}

function operationStarted(state) {
  setLedStatusDoing();

  return {
    ...state,
    failCount: 0,
    activeResultIndex: 0,
    operationStatus: OPERATION_STATUS.DOING
  };
}

function mergeResults(state, data) {
  const rs = state.results;

  if (!data) {
    return rs;
  }

  for (let i = 0; i < data.length; i += 1) {
    rs[i + state.activeResultIndex].ti = data[i].ti;
    rs[i + state.activeResultIndex].mi = data[i].mi;
    rs[i + state.activeResultIndex].wi = data[i].wi;
    rs[i + state.activeResultIndex].result = data[i].result;
    rs[i + state.activeResultIndex].batch = data[i].batch;
  }

  return rs;
}

function operationResultOK(state, action) {
  const { data } = action;
  const results = mergeResults(state, data);

  if (state.operationStatus === OPERATION_STATUS.READY) {
    // 等待车辆状态下 收到结果也无法在进入其他状态
    return {
      ...state,
      activeResultIndex: state.activeResultIndex + data.length,
      failCount: 0,
      // operationStatus: OPERATION_STATUS.DOING,
      results
    };
  }

  return {
    ...state,
    activeResultIndex: state.activeResultIndex + data.length,
    failCount: 0,
    operationStatus: OPERATION_STATUS.DOING,
    results
  };
}

function operationResultNOK(state, action) {
  const { data } = action;
  const results = mergeResults(state, data);

  return {
    ...state,
    failCount: state.failCount + 1,
    results
  };
}

function operationFailed(state, action) {
  const { data } = action;
  const results = mergeResults(state, data);

  if (state.operationStatus === OPERATION_STATUS.READY) {
    // 等待车辆状态下 收到结果也无法在进入其他状态
    return {
      ...state,
      failCount: state.failCount + 1,
      // operationStatus: OPERATION_STATUS.FAIL,
      results
    };
  }
  setLedError(sOn);

  return {
    ...state,
    failCount: state.failCount + 1,
    operationStatus: OPERATION_STATUS.FAIL,
    results
  };
}

function operationFinished(state, action) {
  const { data } = action;

  setLedStatusReady();

  const results = mergeResults(state, data);

  return {
    ...state,
    operationStatus: OPERATION_STATUS.READY,
    carID: '',
    carType: '',
    lnr: '',
    maxOpTimes: 0,
    failCount: 0,
    results
  };
}

function operationContinue(state) {
  if (state.operationStatus === OPERATION_STATUS.READY) {
    // 等待车辆状态下 收到结果也无法在进入其他状态,直接返回
    return state;
  }
  setLedStatusDoing();

  const { activeResultIndex, results } = state;
  let count = 1;
  const ele = results[activeResultIndex + 1];
  for (let i = activeResultIndex + 2; i < results.length; i += 1) {
    if (ele.sequence === results[i].sequence) {
      count += 1;
    } else {
      break;
    }
  }

  return {
    ...state,
    operationStatus: OPERATION_STATUS.DOING,
    activeResultIndex: activeResultIndex + count,
    failCount: 0
  };
}

function switchOperationPredoing(state) {
  return {
    ...state,
    operationStatus: OPERATION_STATUS.PREDOING
  };
}

function switchOperationTimeout(state) {
  return {
    ...state,
    operationStatus: OPERATION_STATUS.TIMEOUT,
    carID: '',
    carType: '',
    lnr: ''
  };
}
