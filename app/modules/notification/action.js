// @flow

import type { CommonLogLvl } from '../../common/utils';

export const NOTIFY = {
  PRE_NEW_NOTIFICATION: 'PRE_NEW_NOTIFICATION',
  NEW_NOTIFICATION: 'NEW_NOTIFICATION',
  CLOSE: 'CLOSE_NOTIFICATION'
};

export type tNotifyVariant = CommonLogLvl;

export const setNewNotification = (variant: tNotifyVariant, message: string, meta) => ({
  type: NOTIFY.NEW_NOTIFICATION,
  variant,
  message,
  meta
});

export const closeNotification = () => ({
  type: NOTIFY.CLOSE
});
