// @flow

export const OPERATION = {
  OPERATION: {
    FETCH_OK: 'OPERATION_OPERATION_FETCH_OK',
    FETCH_FAIL: 'OPERATION_OPERATION_FETCH_FAIL'
  },
  PROGRAMME: {
    SET_OK: 'OPERATION_PROGRAMME_SET_OK',
    SET_FAIL: 'OPERATION_PROGRAMME_SET_FAIL'
  },
  TRIGGER: {
    NEW_DATA: 'OPERATION_TRIGGER_NEW_DATA',
    TRIGGER: 'OPERATION_TRIGGER_TRIGGER',
    BLOCK: 'OPERATION_TRIGGER_BLOCK'
  },
  SOURCE: {
    SET: 'OPERATION_SOURCE_SET'
  },
  RESULT: {
    OK: 'OPERATION_RESULT_OK',
    NOK: 'OPERATION_RESULT_NOK'
  },
  JOB_MANUAL: {
    OK: 'OPERATION_JOB_MANUAL_OK'
  },
  BYPASS: {
    CONFIRM: 'OPERATION_BYPASS_CONFIRM',
    CANCEL: 'OPERATION_BYPASS_CANCEL',
    IO: 'OPERATION_BYPASS_IO'
  },
  CONFLICT: {
    DETECTED: 'OPERATION_CONFLICT_DETECTED',
    CONFIRM: 'OPERATION_CONFLICT_CONFIRM',
    CANCEL: 'OPERATION_CONFLICT_CANCEL'
  },
  STARTED: 'OPERATION_STARTED', // switch to doing
  FINISHED: 'OPERATION_FINISHED', // switch to ready
  FAILED: 'OPERATION_FAILED',
  CONTINUE: 'OPERATION_CONTINUE',
  PREDOING: 'OPERATION_PREDOING',
  TIMEOUT: 'OPERATION_TIMEOUT',
  VERIFIED: 'OPERATION_VERIFIED',
  RESET: 'OPERATION_RESET', // predoing to ready
  SWITCH_STATE: 'OPERATION_SWITCH_STATE'
};

export function switch2Ready(showDiag = true) {
  return {
    type: OPERATION.FINISHED,
    showDiag
  };
}

export function switch2Doing() {
  return {
    type: OPERATION.STARTED
  };
}

// export function operationVerified(data) {
//   return {
//     type: OPERATION.VERIFIED,
//     data
//   };
// }

export function switch2Timeout() {
  return {
    type: OPERATION.TIMEOUT
  };
}

export function switch2PreDoing() {
  return {
    type: OPERATION.PREDOING
  };
}

export function operationTrigger(carID, carType, job, source) {
  return {
    type: OPERATION.TRIGGER.TRIGGER,
    carID,
    carType,
    job,
    source
  };
}

export function operationBypassIO() {
  return {
    type: OPERATION.BYPASS.IO
  };
}

export function operationBypassConfirm() {
  return {
    type: OPERATION.BYPASS.CONFIRM
  };
}

export function operationBypassCancel() {
  return {
    type: OPERATION.BYPASS.CANCEL
  };
}

export function operationConflictDetected(data) {
  return {
    type: OPERATION.CONFLICT.DETECTED,
    data
  };
}

export function operationConflictConfirm(data) {
  return {
    type: OPERATION.CONFLICT.CONFIRM,
    data
  };
}

export function operationConflictCancel() {
  return {
    type: OPERATION.CONFLICT.CANCEL
  };
}

export function operationTriggerBlock(block) {
  return {
    type: OPERATION.TRIGGER.BLOCK,
    block
  };
}
