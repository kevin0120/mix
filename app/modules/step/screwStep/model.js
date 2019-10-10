// @flow

import { find, isNil } from 'lodash-es';
import { CommonLog } from '../../../common/utils';
import { POINT_STATUS, RESULT_STATUS, controllerModes } from './constants';


export type tPointStatus = $Keys<typeof POINT_STATUS>;

export type tResultStatus = $Keys<typeof RESULT_STATUS>;

// NOTE: 拧紧点key定义,如果此点为key，此点必须拧紧结束同时此组的结果到达拧紧关键点个数才能进入下个拧紧组

class ClsOperationPoint {
  _point: tPoint;

  _toolSN: string;

  _isActive: boolean = false;

  _results: Array<tResult>;

  constructor(p: tPoint) {
    this._point = p;
    this._toolSN = p.toolSN;
    this._results = [];
  }

  isFinalFail(): boolean {
    // 结果的长度已经达到最大重试次数，同时最后一条结果为fail
    const rsCount = this._results.length;
    return (
      rsCount >= this._point.maxRetryTimes &&
      this._results[rsCount - 1] === RESULT_STATUS.nok
    );
  }

  isPass(): boolean {
    // 是否需要跳到下一个拧紧点
    const rsCount = this._results.length;
    const lastResult = this._results[rsCount - 1];
    return lastResult === RESULT_STATUS.ak2 || this.isFinalFail();
  }

  get isActive(): boolean {
    return this._isActive;
  }

  set isActive(active: boolean) {
    this._isActive = active;
  }

  get toolSN(): string {
    return this._toolSN;
  }

  Active() {
    this.isActive = true;
  }

  get point(): tPoint {
    return this._point;
  }

  get isKey(): boolean {
    return this._point.is_key;
  }

  appendResult(result: tResult): boolean {
    this._results.push(result);
    return true;
  }
}

class ClsOperationPointGroup {
  _groupSeq: number = 0;

  _operationPoints: Array<ClsOperationPoint> = [];

  constructor(s: number) {
    this._groupSeq = s;
  }

  validatePoint(p: tPoint): boolean {
    const { toolSN } = p;
    const isExist = find(this._operationPoints, { toolSN: toolSN });
    if (isNil(isExist)) {
      return true;
    }
    return false;
  }

  appendNewOperationPoint(p: tPoint): boolean {
    if (!this.validatePoint(p)) {
      CommonLog.lError(
        `appendNewOperationPoint validatePoint Error: ${p.toolSN}`
      );
      return false;
    }
    this._operationPoints.push(new ClsOperationPoint(p));
    this._operationPoints.sort(
      (a: ClsOperationPoint, b: ClsOperationPoint) =>
        a.point.sequence - b.point.sequence
    ); // 排序操作
    return true;
  }

  get groupSequence(): number {
    return this._groupSeq;
  }

  get operationPoints(): Array<ClsOperationPoint> {
    return this._operationPoints;
  }

  set operationPoints(t: Array<ClsOperationPoint>) {
    this._operationPoints = t;
  }

  isGroupPass(): boolean {
    // todo: 是否需要跳到下一个一组拧紧组,
    return true;
  }
}

export class ClsOrderOperationPoints {
  _operationGroups: { [groupSeq: number]: ClsOperationPointGroup } = {};

  static validatePayload(payload: tScrewStepPayload): boolean {
    let ret: boolean = true;
    const { controllerMode, jobID } = payload;
    switch (controllerMode) {
      case 'job':
        if (jobID <= 0) {
          return false;
        }
        break;
      case 'pset':
        // do nothing
        break;
      default:
        CommonLog.lError(`controllerMode: ${controllerMode} Is Not Support!`);
        ret = false;
    }
    return ret;
  }

  constructor(p: tScrewStepPayload) {
    const ret: boolean = ClsOrderOperationPoints.validatePayload(p);
    if (!ret) {
      // 验证失败
      CommonLog.lError(`validatePayload Error! Payload: ${String(p)}`);
      return;
    }
    const { points } = p;
    points.forEach((point: tPoint) => {
      this._appendNewOperationPoint(point);
    });
  }

  _appendNewOperationPoint(p: tPoint) {
    const groupSeq: number = p.group_sequence;
    const isExist = Object.hasOwnProperty.call(this.operationGroups, groupSeq);
    let pg: ClsOperationPointGroup | null = null;
    switch (isExist) {
      case true:
        pg = this.operationGroups[groupSeq];
        break;
      case false:
        CommonLog.Debug(
          `Group Seq: ${groupSeq} Is Not Exist! We Will Create It`
        );
        pg = new ClsOperationPointGroup(groupSeq);
        this._operationGroups[groupSeq] = pg;
        break;
      default:
        CommonLog.lError(
          'tOrderOperationPoints _appendNewOperationPoint Error!'
        );
    }
    if (pg) {
      pg.appendNewOperationPoint(p); // always success
    }
  }

  get operationGroups(): { [groupSeq: number]: ClsOperationPointGroup } {
    return this._operationGroups;
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


export type tControllerMode = $Keys<typeof controllerModes>;

export type tScrewStepPayload = {
  model: string, // 产品类型
  points: Array<tPoint>, // 拧紧点
  jobID: number, // 拧紧ID, 只有当是job模式才会使用
  lnr: string,
  controllerMode: tControllerMode
};
