// @flow

import { NOTIFY } from './actionTypes';

export const setNewNotification = (variant: string, message: string) => ({
  type: NOTIFY.PRE_NEW_NOTIFICATION,
  variant,
  message
});

export const NewNotificationOK = (variant: string, message: string) => ({
  type: NOTIFY.NEW_NOTIFICATION,
  variant,
  message
});

export const closeNotification = () => ({
  type: NOTIFY.CLOSE
});
