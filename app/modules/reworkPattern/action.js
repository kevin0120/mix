// @flow

import { REWORK_PATTERN } from './constants';
import { ClsOperationPoint } from '../step/screwStep/classes/ClsOperationPoint';
import type { IScrewStep } from '../step/screwStep/interface/IScrewStep';


export default {
  tryRework: (order, step: IScrewStep, point: ClsOperationPoint) => ({
    type: REWORK_PATTERN.TRY_REWORK,
    order,
    step,
    point
  }),
  doRework: (order, step, point) => ({
    type: REWORK_PATTERN.DO_REWORK,
    order, step, point
  }),
  cancelRework: () => ({
    type: REWORK_PATTERN.CANCEL_REWORK
  }),
  selectTool: (tool) => ({
    type: REWORK_PATTERN.SELECT_TOOL,
    tool
  })
};
