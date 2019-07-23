import { DIALOG } from './action';

const initState = {
  open: false,
  config: {
    hasOk: false,
    hasCancel: true,
    cancelAction: () => {},
    okAction: () => {}
  }
};

export default function(state = initState, action) {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [DIALOG.SHOW]: (state, action) => ({
    open: true,
    config: action.config
  }),
  [DIALOG.CANCEL]: () => ({}),
  [DIALOG.OK]: () => ({}),
  [DIALOG.CLOSE]: () => ({
    open: false
  })
};
