import { WORKCENTER_MODE } from './constants';
import type { tWorkCenterMode } from './interface/typeDef';
import { translation as trans } from '../../components/NavBar/local';
import type { tCommonActionType } from '../../common/type';


const initState: tWorkCenterMode = trans.normWorkCenterMode;

type tAction = tCommonActionType & { mode: tWorkCenterMode };

export default function(state: tWorkCenterMode = initState, action: tAction): tWorkCenterMode {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [WORKCENTER_MODE.SWITCH]: (state: tWorkCenterMode, action: tAction) => (
    action.mode
  )
};
