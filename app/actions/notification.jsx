// @flow

import {NOTIFY} from "./actionTypes";

export const setNewNotification = (variant: string, message: string) => (
  {
    type: NOTIFY.NEW_NOTIFICATION,
    variant,
    message
  }
);

export const closeNotification = () => (
  {
    type: NOTIFY.NEW_NOTIFICATION,
  }
);
