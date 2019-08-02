import type { CommonLogLvl } from '../../common/utils';

export const NOTIFIER = {
  ENQUEUE_SNACKBAR: 'NOTIFIER_ENQUEUE_SNACKBAR',
  CLOSE_SNACKBAR: 'NOTIFIER_CLOSE_SNACKBAR',
  REMOVE_SNACKBAR: 'NOTIFIER_REMOVE_SNACKBAR'
};

type tnotiStackVariant = 'default' | 'success' | 'error' | 'info';

function convertLogLvl2notiStackVariant(lvl: CommonLogLvl): tnotiStackVariant {
  switch (lvl) {
    case 'Info':
      return 'info';
    case 'Debug':
      return 'error';
    case 'Error':
      return 'error';
    case 'Warn':
      return 'info';
    case 'Maintenance':
      return 'info';
    default:
      return 'default';
  }
}

const enqueueSnackbar = (variant: CommonLogLvl, msg: string) => ({
  type: NOTIFIER.ENQUEUE_SNACKBAR,
  notification: {
    message: msg,
    options: {
      key: `${new Date().getTime() + Math.random()}`,
      variant: convertLogLvl2notiStackVariant(variant),
      anchorOrigin: {
        vertical: 'top',
        horizontal: 'left'
      }
    }
  }
});

const closeSnackbar = key => ({
  type: NOTIFIER.CLOSE_SNACKBAR,
  dismissAll: !key, // dismiss all if no key has been defined
  key
});

const removeSnackbar = key => ({
  type: NOTIFIER.REMOVE_SNACKBAR,
  key
});
export default {
  enqueueSnackbar,
  closeSnackbar,
  removeSnackbar
};
