import { HEALTHZ } from './action';

const initState = {
  status: {}
};

const reducers = {
  [HEALTHZ.DATA]: (state, action) => {
    const { status } = action;
    return {
      status
    };
  }
};

export default function(state = initState, action) {
  return reducers[action.type] && reducers[action.type](state, action) || state;
};