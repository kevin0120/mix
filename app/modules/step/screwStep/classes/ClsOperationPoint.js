// NOTE: 拧紧点key定义,如果此点为key，此点必须拧紧结束同时此组的结果到达拧紧关键点个数才能进入下个拧紧组
import {RESULT_STATUS} from '../constants';
import type { tPoint, tResult } from '../interface/typeDef';

// eslint-disable-next-line import/prefer-default-export
export class ClsOperationPoint {
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