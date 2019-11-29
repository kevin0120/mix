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
    CommonLog.Debug(
      `ClsOrderOperationPoints validatePayload Points: ${JSON.stringify(
        payload.tightening_points
      )}`
    );
    let ret: boolean = true;
    const { controllerMode, jobID } = payload;
    switch (controllerMode) {
      case 'job':
        if (isNil(jobID)) {
          CommonLog.Info(
            'ClsOrderOperationPoints validatePayload Controller Mode Is Job, But Job ID Is Undefined'
          );
          return false;
        }
        if (jobID <= 0) {
          CommonLog.Info(
            'ClsOrderOperationPoints validatePayload Job Is Less Than Zero'
          );
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
      CommonLog.Debug(`Group Seq: ${groupSeq} Does Not Exist! Creating...`);
      pg = new ClsOperationPointGroup(groupSeq, p.key_num);
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
    let newActivePoints: Array<?ClsOperationPoint> = [];
    let newInactivePoints: Array<?ClsOperationPoint> = [];
    results.forEach(r => {
      const { sequence: seq } = r;
      const group = this.getGroupByPointSequence(seq);
      if (!group) {
        return;
      }
      const groupSeq = group.groupSequence;
      const inactivePoints = group.newResult(r);
      newInactivePoints = newInactivePoints.concat(inactivePoints);
      // if (!group.isKeyPass) {
      //   return;
      // }
      // const gSeq = Math.min(
      //   ...Object.keys(this._groups)
      //     .map(s => parseInt(s, 10))
      //     .filter(s => s > groupSeq)
      // );
      // const nextGroup = this._groups[gSeq];
      // if (!nextGroup) {
      //   return;
      // }
      // nextGroup.setActive(true);
      // newActivePoints = newActivePoints.concat(nextGroup.points);
    });

    return {
      // active: newActivePoints,
      inactive: newInactivePoints
    };
  }

  start(): Array<ClsOperationPoint> {
    // 开始可被最先开始的的组，并返回所有被开始的点
    return this.nextActiveGroups().reduce((activatedPoints, g) => {
      const points = g.start();
      console.warn('points', points);
      return activatedPoints.concat(points);
    }, []);
  }

  stop() {
    Object.values(this._groups).forEach(g => {
      g.setActive(false);
    });
  }

  nextActiveGroups(seq) {
    // 下一步可工作的组，即之前所有组都完成的组
    const groups = [];
    const sortedSeqs = Object.keys(this._groups).map(s => parseInt(s, 10))
      .sort((s1, s2) => s1 - s2);
    // eslint-disable-next-line no-restricted-syntax
    for (const s of sortedSeqs) {
      if (seq && s <= seq) {
        // eslint-disable-next-line no-continue
        continue;
      }
      const group = this._groups[s];
      if (!group.isAllPass) {
        groups.push(this._groups[s]);
      }
      if (!group.isKeyPass) {
        break;
      }
    }
    console.warn('groups', groups);
    return groups;
  }


  nextGroupSequence(groupSequence) {
    const sequencesAfter = Object.keys(this._groups)
      .map(s => parseInt(s, 10))
      .filter(s => s > groupSequence);
    return Math.min(...sequencesAfter);
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

  get currentActivePoints(): Array<ClsOperationPoint> {
    return this.points.filter((p: ClsOperationPoint) => p.isActive) || [];
  }

  getGroupByPointSequence(seq: number): ?ClsOperationPointGroup {
    // eslint-disable-next-line flowtype/no-weak-types
    const groups: Array<ClsOperationPointGroup> = (Object.values(
      this.operationGroups
    ): any);
    return groups.find(g => g.points.some(p => p.sequence === seq));
  }

  get isPass() {
    return Object.keys(this._groups).every(g =>
      this._groups[parseInt(g, 10)].isAllPass
    );
  }

  get isFailed() {
    return Object.keys(this._groups).some(g =>
      this._groups[parseInt(g, 10)].isFailed
    );
  }

  hasPoint(point: ClsOperationPoint) {
    return Object.keys(this._groups).some(k =>
      this._groups[parseInt(k, 10)].hasPoint(point)
    );
  }
}
