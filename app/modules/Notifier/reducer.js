// @flow
import { NOTIFIER } from './action';
import type { tNotification, tNotifyState } from './typeDef';
import type { tAction } from '../typeDef';

const initState: tNotifyState = {
  notifications: []
};

const logLvl2NotiStackVariant = {
  Info: 'info',
  Debug: 'error',
  Error: 'error',
  Warn: 'warning',
  Maintenance: 'info',
  default: 'info'
};

export default (state: tNotifyState = initState, action: tAction<any, any>) => {
  switch (action.type) {
    case NOTIFIER.ENQUEUE_SNACKBAR: {
      const { message, variant } = action;
      const notification = {
        message,
        options: {
          key: `${new Date().getTime() + Math.random()}`,
          variant: logLvl2NotiStackVariant[variant],
          anchorOrigin: {
            vertical: 'top',
            horizontal: 'left'
          }
        }
      };
      return {
        ...state,
        notifications: [
          ...state.notifications,
          {
            key: notification.options.key,
            ...notification
          }
        ]
      };
    }
    case NOTIFIER.CLOSE_SNACKBAR:
      return {
        ...state,
        notifications: state.notifications.map<tNotification>(notification =>
          action.dismissAll || notification.key === action.key
            ? { ...notification, dismissed: true }
            : { ...notification }
        )
      };

    case NOTIFIER.REMOVE_SNACKBAR:
      return {
        ...state,
        notifications: state.notifications.filter(
          notification => notification.key !== action.key
        )
      };

    default:
      return state;
  }
};
