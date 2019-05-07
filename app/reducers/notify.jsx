// @flow

import { NOTIFY } from '../actions/actionTypes';

// 当前工作模式、vin/零件号（如果有），车型代码、api及返回值、其他业务相关细信息
const defaultNotify = {
  variant: 'success', // 'error', 'warning', 'info'
  message: '',
  isShow: false,
  meta:{}
};

type actionType = {
  +type: string,
  +variant: string,
  +message: string
};

export default function notify(
  state: object = defaultNotify,
  action: actionType
) {
  switch (action.type) {
    case NOTIFY.NEW_NOTIFICATION: {
      const { type, ...data } = action;
      return {
        ...data,
        isShow: true,
      };
    }
    case NOTIFY.CLOSE: {
      return { ...state, isShow: false };
    }
    default:
      return state;
  }
}
