// @flow

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
  lsn: 'LSN'
};

export type tPointStatus = $Keys<typeof POINT_STATUS>;

export type tResultStatus = $Keys<typeof RESULT_STATUS>;

// NOTE: 拧紧点key定义,如果此点为key，此点必须拧紧结束同时此组的结果到达拧紧关键点个数才能进入下个拧紧组

class tClsOperationPoint {

  _point: tPoint;


  _results: Array<tResult>;

  construct(p: tPoint) {
    this._point = p;
    this._results = [];
  }

  is_final_fail(): boolean {
    const rs_count = this._results.length;
    if(rs_count >= this._point.maxRetryTimes && this._results[rs_count - 1] === RESULT_STATUS.nok) {
      return true;
    }
    return false;
  }

}


export type tPoint = {
  id: number,
  toolSN: string,
  pset: number,
  maxRetryTimes: number, // max_redo_times

  x: number,
  y: number,

  sequence: number,
  group_sequence: number,
  key_num: number, // key_num,定义了关键点个数
  is_key: boolean, // 定义了此点是否为关键点
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
