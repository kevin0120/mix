// @flow
import { isNil } from 'lodash-es';
import { CommonLog } from '../../../../common/utils';
import type { tPoint, tResult, tScrewStepPayload } from '../interface/typeDef';
import { ClsOperationPointGroup } from './ClsOperationPointGroup';
import { ClsOperationPoint } from './ClsOperationPoint';

// eslint-disable-next-line import/prefer-default-export
export class ClsOrderOperationPoints {
  _groups: { [groupSeq: number]: ClsOperationPointGroup } = {};

  static validatePayload(payload: tScrewStepPayload): boolean {
    CommonLog.Debug(`ClsOrderOperationPoints validatePayload Points: ${JSON.stringify(payload.points)}`);
    let ret: boolean = true;
    const { controllerMode, jobID } = payload;
    switch (controllerMode) {
      case 'job':
        if (isNil(jobID)) {
          CommonLog.Info('ClsOrderOperationPoints validatePayload Controller Mode Is Job, But Job ID Is Undefined');
          return false;
        }
        if (jobID <= 0) {
          CommonLog.Info('ClsOrderOperationPoints validatePayload Job Is Less Than Zero');
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

  constructor(points: Array<tPoint>) {
    // const ret: boolean = ClsOrderOperationPoints.validatePayload(p);
    // if (!ret) {
    //   // 验证失败
    //   throw new Error(`validatePayload Error! Payload! `);
    //   // CommonLog.lError(`validatePayload Error! Payload! `);
    //   // return;
    // }
    // const { points } = p;
    points.forEach((point: tPoint) => {
      this._appendNewOperationPoint(point);
    });
  }

  _appendNewOperationPoint(p: tPoint) {
    const groupSeq: number = p.group_sequence;
    const isExist = Object.hasOwnProperty.call(this.operationGroups, groupSeq);
    let pg: ClsOperationPointGroup | null = null;
    if (isExist) {
      pg = this.operationGroups[groupSeq];
    } else {
      CommonLog.Debug(
        `Group Seq: ${groupSeq} Does Not Exist! Creating...`
      );
      pg = new ClsOperationPointGroup(groupSeq,p.minPassCount);
      this._groups[groupSeq] = pg;
    }
    if (pg) {
      pg.appendNewOperationPoint(p); // always success
    }
  }

  get operationGroups(): { [groupSeq: number]: ClsOperationPointGroup } {
    return this._groups;
  }

  newResult(results: Array<tResult>) {
    let newActivePoints = [];
    let newInactivePoints = [];
    results.forEach((r) => {
      const { group_sequence: groupSeq } = r;
      const group = this._groups[groupSeq];
      if (!group) {
        return;
      }
      const point = group.newResult(r);
      newInactivePoints = newInactivePoints.concat(point);
      if (!group.isKeyPass()) {
        return;
      }
      // eslint-disable-next-line radix
      const gSeq = Math.min(...Object.keys(this._groups).map(s => parseInt(s)).filter(s => s > groupSeq));
      const nextGroup = this._groups[gSeq];
      if (!nextGroup) {
        return;
      }
      nextGroup.setActive(true);
      newActivePoints = newActivePoints.concat(nextGroup.points);
    });

    return {
      active: newActivePoints,
      inactive: newInactivePoints
    };
  }

  start(): Array<ClsOperationPoint> {
    // eslint-disable-next-line radix
    const groupSeqs = Object.keys(this._groups).map(s => parseInt(s));
    const firstGroupSeq = Math.min(...groupSeqs);
    this._groups[firstGroupSeq].setActive(true);
    return this._groups[firstGroupSeq].points;
  }

  get points() {
    // eslint-disable-next-line radix
    const groupSeqs = Object.keys(this._groups).map(s => parseInt(s));
    let points = [];
    groupSeqs.forEach(g => {
      points = points.concat(this._groups[g].points);
    });
    return points;
  }

  isPass() {
    // eslint-disable-next-line radix
    return Object.keys(this._groups).every(g => this._groups[parseInt(g)].isAllPass());
  }

  isFailed() {
    // eslint-disable-next-line radix
    return Object.keys(this._groups).some(g => this._groups[parseInt(g)].isFailed());
  }

  hasPoint(point: ClsOperationPoint) {
    // eslint-disable-next-line radix
    return Object.keys(this._groups).some(k => this._groups[parseInt(k)].hasPoint(point));
  }
}
