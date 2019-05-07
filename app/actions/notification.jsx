// @flow

import { NOTIFY } from './actionTypes';

export const setNewNotification = (variant: string, message: string, meta) => ({
  type: NOTIFY.NEW_NOTIFICATION,
  variant,
  message,
  meta
});

export const closeNotification = () => ({
  type: NOTIFY.CLOSE
});
