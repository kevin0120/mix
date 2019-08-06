// @flow
import type { tResult } from './model';

export const SCREW_STEP = {
  RESULT: 'SCREW_STEP_RESULT'
};

export type tResultAction = { type: string, results: Array<tResult> };

export default {
  result: (results: Array<tResult>): tResultAction => ({
    type: SCREW_STEP.RESULT,
    results
  })
};
