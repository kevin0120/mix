// @flow

// NOTE: 拧紧点key定义,如果此点为key，此点必须拧紧结束同时此组的结果到达拧紧关键点个数才能进入下个拧紧组
import { POINT_STATUS, RESULT_STATUS } from '../constants';
import type { tPoint, tPointStatus, tResult } from '../interface/typeDef';

// eslint-disable-next-line import/prefer-default-export
export class ClsOperationPoint {
  _results: Array<tResult>;

  constructor(p: tPoint) {
    this._point = p;
    this._toolSN = p.tightening_tool;
    this._results = [];
  }

  _point: tPoint;

  get point(): tPoint {
    return this._point;
  }

  _toolSN: string;

  get toolSN(): string {
    return this._toolSN;
  }

  _isActive: boolean = false;

  get isActive(): boolean {
    return this._isActive;
  }

  _status: tPointStatus = POINT_STATUS.WAITING;

  get status(): tPointStatus {
    return this._status;
  }

  get isFinalFail(): boolean {
    // 结果的长度已经达到最大重试次数，同时最后一条结果为fail
    // const rsCount = this._results.length;
    // return (
    //   this._point.maxRetryTimes >= 0 &&
    //   rsCount >= this._point.maxRetryTimes &&
    //   this._results[rsCount - 1].result === RESULT_STATUS.nok
    // );
    return (
      this._point.max_redo_times >= 0
      && this._results.filter(r => r.measure_result === RESULT_STATUS.nok)
        .length >= this._point.max_redo_times
      && this._results.slice(-1)[0].measure_result === RESULT_STATUS.nok
    );
  }

  get isPass(): boolean {
    // 是否需要跳到下一个拧紧点
    const rsCount = this._results.length;
    const lastResult: tResult = this._results[rsCount - 1];
    if (!lastResult) {
      return false;
    }
    return (
      lastResult.measure_result === RESULT_STATUS.ak2
      || this.isFinalFail
      || lastResult.measure_result === RESULT_STATUS.ok
    );
  }

  get isSuccess(): boolean {
    const rsCount = this._results.length;
    const lastResult: tResult = this._results[rsCount - 1];
    if (!lastResult) {
      return false;
    }
    return (
      lastResult.measure_result === RESULT_STATUS.ok
    );
  }

  get sequence() {
    return this._point.sequence;
  }

  get noRedo(): boolean {
    return this.status === POINT_STATUS.SUCCESS;
  }

  get x() {
    return this._point.x;
  }

  get y() {
    return this._point.y;
  }

  // eslint-disable-next-line camelcase
  get group_sequence() {
    return this._point.group_sequence;
  }

  get isKey(): boolean {
    return this._point.is_key;
  }

  get data() {
    return {
      ...this.point,
      status: this._status,
      isActive: this._isActive,
      noRedo: this.noRedo
    };
  }

  get pset() {
    return this._point.pset;
  }

  start(forceStart) {
    if (this.isActive) {
      return null;
    }
    if (forceStart || !this.isPass) {
      this.setActive(true);
      return this;
    }
    return null;
  }

  toString(): string {
    return JSON.stringify(
      {
        Bolt: this.point.nut_no,
        Sequence: this.point.sequence
      });
  }

  setActive(active: boolean) {
    this._isActive = active;
  }

  newResult(result: tResult): ?ClsOperationPoint {
    const r = {
      ...result,
      measure_result: result.measure_result.toUpperCase()
    };
    this._results.push(r);

    this._parseStatus(r);
    this._parseActive(r);
    if (!this._isActive) {
      return this;
    }
    return null;
  }

  _parseStatus(result: tResult): void {
    if (!result.measure_result || typeof result.measure_result !== 'string') {
      return;
    }

    if (result.measure_result === RESULT_STATUS.ok) {
      this._status = POINT_STATUS.SUCCESS;
    }
    if (result.measure_result === RESULT_STATUS.nok) {
      this._status = POINT_STATUS.ERROR;
    }
  }

  _parseActive(result: tResult): void {
    if (this.isPass) {
      this.setActive(false);
    }
  }

  redo() {
    this.setActive(true);
    this._results = [];
    return this;
  }
}
