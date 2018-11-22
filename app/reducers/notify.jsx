// @flow

import { NOTIFY } from '../actions/actionTypes';


const defaultNotify = {
  variant: 'success', // 'error', 'warning', 'info'
  message: '',
  isShow: false
};

type actionType = {
  +type: string,
  +variant: string,
  +message: string,
};

export default function notify(
  state: object = defaultNotify,
  action: actionType
) {
  switch (action.type) {
    case NOTIFY.NEW_NOTIFICATION:{
      const {variant, message} = action;
      return {variant, message, isShow: true};
    }
    case NOTIFY.CLOSE: {
      return {...state, isShow: false};
    }
    default:
      return state;
  }
}
