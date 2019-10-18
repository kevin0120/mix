// @flow
import type { CommonLogLvl } from '../../common/utils';
import type { tNotifyKey } from './typeDef';

export type tNotifyVariant = CommonLogLvl;

export const NOTIFIER = {
  ENQUEUE_SNACKBAR: 'NOTIFIER_ENQUEUE_SNACKBAR',
  CLOSE_SNACKBAR: 'NOTIFIER_CLOSE_SNACKBAR',
  REMOVE_SNACKBAR: 'NOTIFIER_REMOVE_SNACKBAR'
};

const enqueueSnackbar = (
  variant: CommonLogLvl,
  message: string,
  meta: Object
) => ({
  type: NOTIFIER.ENQUEUE_SNACKBAR,
  meta,
  variant,
  message
});

const closeSnackbar = (key: tNotifyKey) => ({
  type: NOTIFIER.CLOSE_SNACKBAR,
  dismissAll: !key, // dismiss all if no key has been defined
  key
});

const removeSnackbar = (key: string) => ({
  type: NOTIFIER.REMOVE_SNACKBAR,
  key
});

export default {
  enqueueSnackbar,
  closeSnackbar,
  removeSnackbar
};
