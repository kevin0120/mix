// @flow
import type { tResult } from './model';

export const SCREW_STEP = {
  RESULT: 'SCREW_STEP_RESULT',
  REDO_POINT:'SCREW_STEP_REDO_POINT'
};

export type tResultAction = { type: string, results: Array<tResult> };

export default {
  result: (results: Array<tResult>): tResultAction => ({
    type: SCREW_STEP.RESULT,
    results
  }),
  redoPoint:(point): tResultAction => ({
    type: SCREW_STEP.REDO_POINT,
    point
  }),
};
