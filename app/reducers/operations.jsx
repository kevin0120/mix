import { OPERATION } from '../actions/actionTypes';
import { isCarID } from '../common/utils';

import {
  setLedStatusDoing,
  setLedError,
  setLedStatusReady,
  sOn,
  sBlinkOn
} from '../sagas/io';

export const OPERATION_STATUS = {
  INIT: 'Init',
  READY: 'Ready',
  PREDOING: 'PreDoing',
  DOING: 'Doing',
  TIMEOUT: 'Timeout',
  FAIL: 'Fail'
};

export const OPERATION_RESULT = {
  OK: 'OK',
  NOK: 'NOK'
};

export const OPERATION_SOURCE = {
  SCANNER: 'SCANNER',
  RFID: 'RFID',
  ANDON: 'ANDON',
  MANUAL: 'MANUAL'
};

const defaultOperations = {
  operationStatus: 'Ready',
  carID: '',
  carType: '',
  activeResultIndex: 0,
  failCount: 0,
  jobID: 0,
  maxOpTimes: 30,
  workSheet: '',
  productID: -1,
  workcenterID: -1,
  lnr: '',
  source: '',
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
    //   result: ''
    // }
  ]
};

type actionType = {
  +type: string,
  +data: object,
  +carID: string,
  +carType: string
};

export default function operations(
  state: object = defaultOperations,
  action: actionType
) {
  switch (action.type) {
    case OPERATION.TRIGGER.NEW_DATA:
      return NewTriggerData(state, action.carID, action.carType);
    case OPERATION.OPERATION.FETCH_OK:
      return NewOperation(state, action.mode, action.data);
    case OPERATION.STARTED:
      return OperationStarted(state);
    case OPERATION.RESULT.OK:
      return OperationResultOK(state, action.data);
    case OPERATION.RESULT.NOK:
      return OperationResultNOK(state, action.data);
    case OPERATION.FAILED:
      return OperationFailed(state, action.data);
    case OPERATION.FINISHED:
      return OperationFinished(state, action.data);
    case OPERATION.CONTINUE:
      return OperationContinue(state);
    case OPERATION.PREDOING:
      return OperationPreDoing(state);
    default:
      return state;
  }
}

function NewTriggerData(state, carID, carType) {
  return {
    ...state,
    carID: carID !== null ? carID : state.carID,
    carType: carType !== null ? carType : state.carType
  };
}

function NewOperation(state, mode, data) {
  if (mode === 'op') {
    // 作业模式
    return {
      ...state,
      jobID: data.job,
      maxOpTimes: data.max_op_time,
      workSheet: data.img,
      productID: data.product_id,
      workcenterID: data.workcenter_id,
      results: data.points
      // activeResultIndex: 0,
      // failCount: 0
    };
  }

  // 工单模式
  return {
    ...state,
    jobID: data.job_id,
    carType: data.model,
    maxOpTimes: data.max_op_time,
    workSheet: data.work_sheet,
    results: data.results,
    // activeResultIndex: 0,
    // failCount: 0,
    lnr: data.lnr
  };
}

function OperationStarted(state) {
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

  for (let i = 0; i < data.length; i++) {
    rs[i + state.activeResultIndex].ti = data[i].ti;
    rs[i + state.activeResultIndex].mi = data[i].mi;
    rs[i + state.activeResultIndex].wi = data[i].wi;
    rs[i + state.activeResultIndex].result = data[i].result;
  }

  return rs;
}

function OperationResultOK(state, data) {
  const results = mergeResults(state, data);

  return {
    ...state,
    activeResultIndex: state.activeResultIndex + data.length,
    failCount: 0,
    operationStatus: OPERATION_STATUS.DOING,
    results
  };
}

function OperationResultNOK(state, data) {
  const results = mergeResults(state, data);

  return {
    ...state,
    failCount: state.failCount + 1,
    results
  };
}

function OperationFailed(state, data) {
  setLedError(sOn);
  const results = mergeResults(state, data);

  return {
    ...state,
    failCount: state.failCount + 1,
    operationStatus: OPERATION_STATUS.FAIL,
    results
  };
}

function OperationFinished(state, data) {
  setLedStatusReady();

  const results = mergeResults(state, data);

  return {
    ...state,
    operationStatus: OPERATION_STATUS.READY,
    results
  };
}

function OperationContinue(state) {
  setLedStatusDoing();

  const { activeResultIndex, results } = state;
  let count = 1;
  const ele = results[activeResultIndex + 1];
  for (let i = activeResultIndex + 2; i < results.length; i++) {
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

function OperationPreDoing(state) {
  return {
    ...state,
    operationStatus: OPERATION_STATUS.PREDOING
  };
}
