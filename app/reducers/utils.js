export function genReducers(reducers, initState={}) {
  return (state = initState, action) => {
    if (reducers[action.type]) {
      return reducers[action.type](state, action);
    }
    return state;
  };
}
