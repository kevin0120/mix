import type { tDialogConfig } from './interface/typeDef';
import {DIALOG} from './constants';

export default {
  dialogShow: (config: tDialogConfig) => ({
    type: DIALOG.SHOW,
    config
  }),
  dialogButton: idx => ({
    type: DIALOG.BUTTON,
    idx
  }),
  dialogClose: () => ({
    type: DIALOG.CLOSE
  })
};
