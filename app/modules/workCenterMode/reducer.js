import { WORKCENTER_MODE, workModes } from './constants';
import type { tWorkCenterMode, tAction } from './interface/typeDef';


const initState: tWorkCenterMode = workModes.normWorkCenterMode;


export default function(state: tWorkCenterMode = initState, action: tAction): tWorkCenterMode {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [WORKCENTER_MODE.SWITCH_VALID]: (state: tWorkCenterMode, action: tAction) => (
    action.mode
  )
};
