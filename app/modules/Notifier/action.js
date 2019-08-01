export const NOTIFIER = {
  ENQUEUE_SNACKBAR: 'NOTIFIER_ENQUEUE_SNACKBAR',
  CLOSE_SNACKBAR: 'NOTIFIER_CLOSE_SNACKBAR',
  REMOVE_SNACKBAR: 'NOTIFIER_REMOVE_SNACKBAR'
};

const enqueueSnackbar = (variant,msg) => {
  return {
    type: NOTIFIER.ENQUEUE_SNACKBAR,
    notification: {
      message: msg,
      options:{
        key: `${new Date().getTime() + Math.random()}`,
        variant,
        anchorOrigin: {
          vertical: 'top',
          horizontal: 'left',
        },
      },
    }
  };
};

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
