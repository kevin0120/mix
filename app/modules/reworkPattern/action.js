// @flow

import { REWORK_PATTERN } from './constants';
import { ClsOperationPoint } from '../step/screwStep/classes/ClsOperationPoint';
import type { IScrewStep } from '../step/screwStep/interface/IScrewStep';
import type { IWorkStep } from '../step/interface/IWorkStep';


export default {
  aReworkSpecialScrewPoint: (step: IWorkStep, point: ClsOperationPoint) => ({
    type: REWORK_PATTERN.SPEC_SCREW,
    extra: { step, point }
  }),
  aReworkScrewStep: (step: IScrewStep) => ({
    type: REWORK_PATTERN.SCREW_STEP,
    extra: step
  })
};
