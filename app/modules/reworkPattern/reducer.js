import { REWORK_PATTERN } from './constants';
import type {  tAction } from './interface/typeDef';

const initState = {
  selectPoints: []
};


export default function(state = initState, action: tAction) {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [REWORK_PATTERN.TRY_REWORK]: (state, action: tAction) => {
    console.warn(action);
    if (!action.point) {
      return state;
    }
    return {
      ...state,
      selectPoints: [action.point]
    };
  },
  [REWORK_PATTERN.CANCEL_REWORK]:(state: tWorkCenterMode, action: tAction)=>({
    ...state,
    selectPoints: []
  })
};
