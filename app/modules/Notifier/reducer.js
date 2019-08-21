import { NOTIFIER } from './action';

const initState = {
  notifications: []
};

export default (state = initState, action) => {
  switch (action.type) {
    case NOTIFIER.ENQUEUE_SNACKBAR:
      return {
        ...state,
        notifications: [
          ...state.notifications,
          {
            key: action.notification.options.key,
            ...action.notification
          }
        ]
      };

    case NOTIFIER.CLOSE_SNACKBAR:
      return {
        ...state,
        notifications: state.notifications.map(notification => {
          return (
            (action.dismissAll || notification.key === action.key)
              ? { ...notification, dismissed: true }
              : { ...notification }
          );
        })
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
