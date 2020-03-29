// @flow
import { isNil } from 'lodash-es';
import { CommonLog } from '../../../../common/utils';
import type { tControl, tPoint, tResult, tScrewStepPayload } from '../interface/typeDef';
import { ClsOperationPointGroup } from './ClsOperationPointGroup';
import { ClsOperationPoint } from './ClsOperationPoint';

// eslint-disable-next-line import/prefer-default-export
export class ClsOrderOperationPoints {
  _groups: { [groupSeq: number]: ClsOperationPointGroup } = {};

  constructor(points: Array<tPoint>) {
    // const ret: boolean = ClsOrderOperationPoints.validatePayload(p);
    // if (!ret) {
    //   // 验证失败
    //   throw new Error(`validatePayload Error! Payload! `);
    //   // CommonLog.lError(`validatePayload Error! Payload! `);
    //   // return;
    // }
    // const { points } = p;
    points.sort((p1, p2) => p1.sequence - p2.sequence).forEach((point: tPoint) => {
      this._appendNewOperationPoint(point);
    });
  }

  get activeControls() {
    return this.points.filter(p => p.isActive).reduce((controls, p) => {
      return controls.concat(p.controls);
    }, []);
  }

  get operationGroups(): { [groupSeq: number]: ClsOperationPointGroup } {
    return this._groups;
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

  newResult(results: Array<tResult>) {
    let newInactiveControls: Array<?ClsOperationPoint> = [];
    results.forEach(r => {
      const { tool_sn, seq } = r;
      const controls = this.activeControls.filter(c => c.toolSN === tool_sn);

      const pointSeqs = [...new Set(controls.map(c => c.sequence))].sort();

      if (!pointSeqs.length > 0) {
        const point = this.points.find(p => p.sequence === seq);
        point.newResult(r);
        return;
      }
      const firstSeq = pointSeqs[0];
      const point = this.points.find(p => p.sequence === firstSeq);
      const inactive = point.newResult(r);
      if (inactive) {
        pointSeqs.shift();
      }
      if (pointSeqs.length > 0) {
        return;
      }
      newInactiveControls = newInactiveControls.concat(...controls);
    });

    return newInactiveControls;
  }

  start(): Array<tControl> {
    // 所有可被最先开始的的组内的点
    const allPoints = this.nextActiveGroups().reduce((points, g) => {
      return points.concat(g.points);
    }, []).filter(p => !p.isPass);


    let activeControls = [];
    this.points.filter(p => p.isActive).forEach(p => {
      activeControls = activeControls.concat(p.controls);
    });

    // 工具未被占用的点
    return [].concat(...allPoints.filter(p => {
      // todo : 有一把工具占用就无法拧紧这个点？
      if (activeControls.find(c => {
        if (p.point.tightening_tools && p.point.tightening_tools.length > 0) {
          return p.toolSNs.find(t => c.toolSN === t) && c.controllerModeId !== p.pset;
        }
        return p.toolSNs.find(t => c.toolSN === t);
      })) {
        return false;
      }
      activeControls = activeControls.concat(p.controls);
      return true;
    }).map(p => p.start())).filter(c => !!c);
  }

  stop() {
    Object.values(this._groups).forEach(g => {
      g.setActive(false);
    });
  }

  byPassControls(controls) {
    controls.forEach(c => {
      const point = this.points.find(pp => pp.sequence === c.sequence);
      if (point) {
        point.setBypass(true);
      }
    });
  }

  clearByPass() {
    this.points.forEach(p => p.setBypass(false));
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
    return groups;
  }

  nextGroupSequence(groupSequence) {
    const sequencesAfter = Object.keys(this._groups)
      .map(s => parseInt(s, 10))
      .filter(s => s > groupSequence);
    return Math.min(...sequencesAfter);
  }

  getGroupByPointSequence(seq: number): ?ClsOperationPointGroup {
    // eslint-disable-next-line flowtype/no-weak-types
    const groups: Array<ClsOperationPointGroup> = (Object.values(
      this.operationGroups
    ): any);
    return groups.find(g => g.points.some(p => p.sequence === seq));
  }

  hasPoint(point: ClsOperationPoint) {
    return Object.keys(this._groups).some(k =>
      this._groups[parseInt(k, 10)].hasPoint(point)
    );
  }
}
