

export const POINT_STATUS = {
  WAITING: 'WAITING',
  WAITING_ACTIVE: 'WAITING_ACTIVE',
  SUCCESS: 'SUCCESS',
  ERROR: 'ERROR',
  ERROR_ACTIVE: 'ERROR_ACTIVE'
};

export const RESULT_STATUS = {
  ok: 'OK',
  nok: 'NOK',
  lsn: 'LSN', // 松脱
  ak2: 'AK2' // 强制放行
};

export const controllerModes = {
  job: 'job',
  pset: 'pset'
};

// screw step action types
export const SCREW_STEP = {
  RESULT: 'SCREW_STEP_RESULT',
  REDO_POINT: 'SCREW_STEP_REDO_POINT',
  CONFIRM_FAIL: 'SCREW_STEP_CONFIRM_FAIL'
};
