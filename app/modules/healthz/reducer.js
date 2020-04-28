// @flow
import { HEALTHZ } from './action';
import type { tHealthzState } from './typeDef';
import type { tAction } from '../typeDef';

const initState = {
  status: {
    rush: false
  }
};

const reducers = {
  [HEALTHZ.DATA]: (state, action) => {
    const { status } = action;
    return {
      status: {
        ...state.status,
        ...status
      }
    };
  }
};

export default function(state: tHealthzState = initState, action: tAction<any, any>) {
  return (
    (reducers[action.type] && reducers[action.type](state, action)) || state
  );
}
