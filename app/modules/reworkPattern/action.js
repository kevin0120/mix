// @flow

import { REWORK_PATTERN } from './constants';
import { ClsOperationPoint } from '../step/screwStep/classes/ClsOperationPoint';
import type { IScrewStep } from '../step/screwStep/interface/IScrewStep';


export default {
  aReworkSpecialScrewPoint: (point: ClsOperationPoint) => ({
    type: REWORK_PATTERN.SPEC_SCREW,
    extra: point
  }),
  aReworkScrewStep: (step: IScrewStep) => ({
    type: REWORK_PATTERN.SCREW_STEP,
    extra: step
  }),
};
