// @flow

// NOTE: 拧紧点key定义,如果此点为key，此点必须拧紧结束同时此组的结果到达拧紧关键点个数才能进入下个拧紧组
import { POINT_STATUS, RESULT_STATUS } from '../constants';
import type { tPoint, tPointStatus, tResult } from '../interface/typeDef';

// eslint-disable-next-line import/prefer-default-export
export class ClsOperationPoint {
  _point: tPoint;

  _toolSN: string;

  _isActive: boolean = false;

  _status: tPointStatus = POINT_STATUS.WAITING;

  _results: Array<tResult>;

  constructor(p: tPoint) {
    this._point = p;
    this._toolSN = p.tightening_tool;
    this._results = [];
  }

  isFinalFail(): boolean {
    // 结果的长度已经达到最大重试次数，同时最后一条结果为fail
    // const rsCount = this._results.length;
    // return (
    //   this._point.maxRetryTimes >= 0 &&
    //   rsCount >= this._point.maxRetryTimes &&
    //   this._results[rsCount - 1].result === RESULT_STATUS.nok
    // );
    return (
      this._point.max_redo_times >= 0 &&
      this._results.filter(r => r.result === RESULT_STATUS.nok).length >=
        this._point.max_redo_times
    );
  }

  _isPass(): boolean {
    // 是否需要跳到下一个拧紧点
    const rsCount = this._results.length;
    const lastResult = this._results[rsCount - 1];
    return lastResult === RESULT_STATUS.ak2 || this.isFinalFail();
  }

  get sequence() {
    return this._point.sequence;
  }

  get status(): tPointStatus {
    return this._status;
  }

  get isActive(): boolean {
    return this._isActive;
  }
  
  get canRedo(): boolean {
    return !this.isActive && this.status === POINT_STATUS.ERROR;
  }

  get toolSN(): string {
    return this._toolSN;
  }

  get x() {
    return this._point.x;
  }

  get y() {
    return this._point.y;
  }

  get group_sequence() {
    return this._point.group_sequence;
  }

  setActive(active: boolean) {
    this._isActive = active;
  }

  get point(): tPoint {
    return this._point;
  }

  get isKey(): boolean {
    return this._point.is_key;
  }

  newResult(result: tResult): ?ClsOperationPoint {
    this._results.push(result);

    this._parseStatus(result);
    this._parseActive(result);
    if (!this._isActive) {
      return this;
    }
    return null;
  }

  _parseStatus(result: tResult): void {
    if(!result.result || typeof result.result !=='string'){
      return;
    }

    if (result.result.toUpperCase() === RESULT_STATUS.ok.toUpperCase()) {
      this._status = POINT_STATUS.SUCCESS;
    }
    if (result.result.toUpperCase() === RESULT_STATUS.nok.toUpperCase()) {
      this._status = POINT_STATUS.ERROR;
    }
  }

  _parseActive(result: tResult): void {
    if (this.isFinalFail() || result.result === RESULT_STATUS.ok) {
      this.setActive(false);
    }
  }

  redo() {
    this.setActive(true);
    this._results = [];
    return this;
  }
}
