import type { dialogConfig } from './model';

export const DIALOG = {
  SHOW: 'DIALOG_SHOW',
  CANCEL: 'DIALOG_CANCEL',
  OK: 'DIALOG_OK',
  CLOSE: 'DIALOG_CLOSE'
};

export default {
  showDialog: (config: dialogConfig) => ({
    type: DIALOG.SHOW,
    config
  }),
  cancelDialog: () => ({
    type: DIALOG.CANCEL
  }),
  okDialog: () => ({
    type: DIALOG.OK
  }),
  closeDialog: () => ({
    type: DIALOG.CLOSE
  })
};
