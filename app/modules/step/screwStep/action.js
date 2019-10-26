// @flow
import type {tResult } from './interface/typeDef';
import { SCREW_STEP } from './constants';
import { ClsOperationPoint } from './classes/ClsOperationPoint';

export type tResultAction = { type: string, results: Array<tResult> };

export default {
  result: (results: Array<tResult>): tResultAction => ({
    type: SCREW_STEP.RESULT,
    results
  }),
  redoPoint: (point: ClsOperationPoint) => ({
    type: SCREW_STEP.REDO_POINT,
    point
  }),
  confirmFail: () => ({
    type: SCREW_STEP.CONFIRM_FAIL
  })
};
