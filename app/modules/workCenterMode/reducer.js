import { WORKCENTER_MODE } from './constants';
import type { tWorkCenterMode, tAction } from './interface/typeDef';
import { translation as trans } from '../../components/NavBar/local';


const initState: tWorkCenterMode = trans.normWorkCenterMode;


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
