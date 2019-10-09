import { HEALTHZ } from './action';

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

export default function(state = initState, action) {
  return (
    (reducers[action.type] && reducers[action.type](state, action)) || state
  );
}
