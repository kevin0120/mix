import type { tWorkCenterMode } from './interface/typeDef';
import { WORKCENTER_MODE } from './constants';

export default {
  aSwitchWorkCenterMode: (mode: tWorkCenterMode) => ({
    type: WORKCENTER_MODE.SWITCH,
    mode
  }),
  aSwitchWorkCenterModeValidOK: (mode: tWorkCenterMode) => ({
    type: WORKCENTER_MODE.SWITCH_VALID,
    mode
  }),
};
