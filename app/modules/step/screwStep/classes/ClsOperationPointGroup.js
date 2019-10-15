import { find, isNil } from 'lodash-es';
import { CommonLog } from '../../../../common/utils';
import type { tPoint } from '../interface/typeDef';
import {ClsOperationPoint} from './ClsOperationPoint';

// eslint-disable-next-line import/prefer-default-export
export class ClsOperationPointGroup {
  _groupSeq: number = 0;

  _operationPoints: Array<ClsOperationPoint> = [];

  constructor(s: number) {
    this._groupSeq = s;
  }

  validatePoint(p: tPoint): boolean {
    const { toolSN } = p;
    const isExist = find(this._operationPoints, { toolSN });
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