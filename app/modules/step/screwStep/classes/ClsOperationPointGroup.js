// @flow
import { find, isNil } from 'lodash-es';
import { CommonLog } from '../../../../common/utils';
import type { tPoint, tResult } from '../interface/typeDef';
import { ClsOperationPoint } from './ClsOperationPoint';
import { POINT_STATUS } from '../constants';

// eslint-disable-next-line import/prefer-default-export
export class ClsOperationPointGroup {
  _groupSeq: number = 0;

  _points: Array<ClsOperationPoint> = [];

  _active: boolean = false;

  _keyNum: number = 0;

  constructor(s: number, keyNum: number) {
    this._groupSeq = s;
    this._keyNum = keyNum;
  }

  validatePoint(p: tPoint): boolean {
    // eslint-disable-next-line camelcase
    const { tightening_tool } = p;
    const pointWithSameToolSN = find(this._points, { tightening_tool });
    return isNil(pointWithSameToolSN);
  }

  appendNewOperationPoint(p: tPoint): boolean {
    if (!this.validatePoint(p)) {
      CommonLog.lError(
        `appendNewOperationPoint validatePoint Error: ${p.tightening_tool}`
      );
      return false;
    }
    this._points.push(new ClsOperationPoint(p));
    // this._points.sort(
    //   (a: ClsOperationPoint, b: ClsOperationPoint) =>
    //     a.point.sequence - b.point.sequence
    // ); // 排序操作
    return true;
  }

  get groupSequence(): number {
    return this._groupSeq;
  }

  get operationPoints(): Array<ClsOperationPoint> {
    return this._points;
  }

  set operationPoints(t: Array<ClsOperationPoint>) {
    this._points = t;
  }

  get points() {
    return this._points;
  }

  // 关键点全部完成
  get isKeyPass() {
    return (
      this._points
        .filter(p => p.isKey)
        .every(p => p.isPass) &&
      this._points.filter(p => p.isPass).length >=
      this._keyNum
    );
  }

  // 所有点完成
  get isAllPass() {
    return this.points.every(p => p.isPass);
  }

  get isFailed() {
    return this.points.some(p => p.isFinalFail);
  }

  start() {
    // 返回所有被开始的点
    return this._points.reduce((activedPoints, p) => {
      const point = p.start();
      if (point) {
        activedPoints.push(point);
      }
      return activedPoints;
    }, []);
  }

  newResult(result: tResult): Array<?ClsOperationPoint> {
    const inactivePoints = [];
    (() => {
      console.log(result);
      const { seq } = result;
      const point = this._points.find(p => p.sequence === seq);
      if (!point) {
        return;
      }
      inactivePoints.push(point.newResult(result));
      if (!this.isAllPass) {
        return;
      }
      this.setActive(false);
    })();
    return inactivePoints.filter(p => !isNil(p));
  }

  setActive(active: boolean) {
    this._active = active;
    this._points.forEach(p => p.setActive(active));
  }

  hasPoint(point: ClsOperationPoint) {
    return this._points.some(p => p === point);
  }
}
