// @flow

import { POINT_STATUS, controllerModes } from '../constants';


export type tPointStatus = $Keys<typeof POINT_STATUS>;

export type tResultStatus = 'OK' | 'NOK' | 'LSN' | 'AK2';

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
  jobID: number,
  timeLine: Array<{
    title: string,
    color: string,
    footerTitle: string,
    body: string
  }>
};


export type tControllerMode = $Keys<typeof controllerModes>;

export type tScrewStepPayload = {
  model: string, // 产品类型
  points: Array<tPoint>, // 拧紧点
  jobID: number, // 拧紧ID, 只有当是job模式才会使用
  lnr: string,
  controllerMode: tControllerMode
};
