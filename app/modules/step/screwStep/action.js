// @flow
import type { tResult } from './model';
import { SCREW_STEP } from './constents';

export type tResultAction = { type: string, results: Array<tResult> };

export default {
  result: (results: Array<tResult>): tResultAction => ({
    type: SCREW_STEP.RESULT,
    results
  }),
  redoPoint: (point): tResultAction => ({
    type: SCREW_STEP.REDO_POINT,
    point
  }),
  confirmFail: () => ({
    type: SCREW_STEP.CONFIRM_FAIL
  })
};
