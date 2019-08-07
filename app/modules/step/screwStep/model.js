// @flow

export const POINT_STATUS = {
  WAITING: 'WAITING',
  WAITING_ACTIVE: 'WAITING_ACTIVE',
  SUCCESS: 'SUCCESS',
  ERROR: 'ERROR',
  ERROR_ACTIVE: 'ERROR_ACTIVE'
};

export const RESULT_STATUS = {
  ok: 'ok',
  nok: 'nok',
  lsn: 'lsn'
};

export type tPointStatus = $Keys<typeof POINT_STATUS>;

export type tResultStatus = $Keys<typeof RESULT_STATUS>;

export type tPoint = {
  id: number,
  toolSN: string,
  pset: number,
  maxRetryTimes: number, // max_redo_times

  x: number,
  y: number,

  sequence: number,
  group_sequence: number,

  ti: number,
  mi: number,
  wi: number,
  status: tPointStatus, // result
  batch: string
};

export type tResult = {
  toolSN: string,
  sequence: number,
  group_sequence: number,
  ti: number,
  mi: number,
  wi: number,
  result: tResultStatus, // result
  batch: string
};

export type tScrewStepData = {
  points: Array<tPoint>,
  retryTimes: number,
  activeIndex: number,
  jobID: number
};

export const controllerModes = {
  job: 'job',
  pset: 'pset'
};

export type tControllerMode = $Keys<typeof controllerModes>;

export type tScrewStepPayload = {
  model: string,
  points: Array<tPoint>,
  jobID: number,
  lnr: string,
  controllerMode: tControllerMode
};
