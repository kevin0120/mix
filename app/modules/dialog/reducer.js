import { DIALOG } from './action';
import type { tDialogState } from './model';

const initState: tDialogState = {
  open: false,
};

export default function(state: tDialogState = initState, action): tDialogState {
  if (reducers[(action?.type)]) {
    return reducers[(action?.type)](state, action);
  }
  return state;
}

const reducers = {
  [DIALOG.SHOW]: (state: tDialogState, action) => ({
    open: true,
    config: action.config
  }),
  [DIALOG.BUTTON]: () => ({}),
  [DIALOG.CLOSE]: () => ({
    open: false
  })
};
