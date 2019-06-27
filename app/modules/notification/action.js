// @flow

export const NOTIFY = {
  PRE_NEW_NOTIFICATION: 'PRE_NEW_NOTIFICATION',
  NEW_NOTIFICATION: 'NEW_NOTIFICATION',
  CLOSE: 'CLOSE_NOTIFICATION'
};

export const setNewNotification = (variant: string, message: string, meta) => ({
  type: NOTIFY.NEW_NOTIFICATION,
  variant,
  message,
  meta
});

export const closeNotification = () => ({
  type: NOTIFY.CLOSE
});
