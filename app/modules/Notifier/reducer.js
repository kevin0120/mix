import { NOTIFIER } from './action';
import type { CommonLogLvl } from '../../common/utils';

const initState = {
  notifications: []
};

type tNotiStackVariant = 'default' | 'success' | 'error' | 'info';

function convertLogLvl2NotiStackVariant(lvl: CommonLogLvl): tNotiStackVariant {
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

export default (state = initState, action) => {
  switch (action.type) {
    case NOTIFIER.ENQUEUE_SNACKBAR: {
      const { message, variant } = action;
      const notification = {
        message,
        options: {
          key: `${new Date().getTime() + Math.random()}`,
          variant: convertLogLvl2NotiStackVariant(variant),
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
        notifications: state.notifications.map(notification =>
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
