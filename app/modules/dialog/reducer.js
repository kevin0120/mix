import { DIALOG } from './action';
import type { dialogState } from './model';

const initState: dialogState = {
  open: false,
  config: {
    hasOk: false,
    hasCancel: true,
    cancelAction: () => {
    },
    okAction: () => {
    }
  }
};

export default function(state: dialogState = initState, action): dialogState {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [DIALOG.SHOW]: (state: dialogState, action) => ({
    open: true,
    config: action.config
  }),
  [DIALOG.CANCEL]: () => ({}),
  [DIALOG.OK]: () => ({}),
  [DIALOG.CLOSE]: () => ({
    open: false
  })
};
