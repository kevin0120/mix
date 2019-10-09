import type { tDialogConfig } from './model';

export const DIALOG = {
  SHOW: 'DIALOG_SHOW',
  BUTTON: 'DIALOG_BUTTON',
  CLOSE: 'DIALOG_CLOSE'
};

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
