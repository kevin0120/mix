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

const reducers = {
  [NOTIFIER.ENQUEUE_SNACKBAR]: (state, action) => {
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
  },
  [NOTIFIER.CLOSE_SNACKBAR]: (state, action) => ({
    ...state,
    notifications: state.notifications.map<tNotification>(notification =>
      action.dismissAll || notification.key === action.key
        ? { ...notification, dismissed: true }
        : { ...notification }
    )
  }),
  [NOTIFIER.REMOVE_SNACKBAR]: (state, action) => ({
    ...state,
    notifications: state.notifications.filter(
      notification => notification.key !== action.key
    )
  })
};

// eslint-disable-next-line flowtype/no-weak-types
export default (state: tNotifyState = initState, action: tAction<any, any>) => {
  if (!reducers[action.type]) {
    return state;
  }
  return reducers[action.type](state, action);
};
