import { isNil } from 'lodash-es';
import { CommonLog } from '../../../../common/utils';
import type { tPoint, tScrewStepPayload } from '../interface/typeDef';
import { ClsOperationPointGroup } from './ClsOperationPointGroup';

// eslint-disable-next-line import/prefer-default-export
export class ClsOrderOperationPoints {
  _operationGroups: { [groupSeq: number]: ClsOperationPointGroup } = {};

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

  constructor(p: tScrewStepPayload) {
    const ret: boolean = ClsOrderOperationPoints.validatePayload(p);
    if (!ret) {
      // 验证失败
      CommonLog.lError(`validatePayload Error! Payload! `);
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
