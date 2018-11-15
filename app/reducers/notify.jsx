// @flow

const defaultNotify = {
  variant: 'success', // 'error', 'warning', 'info'
  message: '',
  isShow: false
};

type actionType = {
  +type: string,
  +variant: string,
  +message: string,
  +isShow: boolean
};

export default function notify(state: object = defaultNotify, action: actionType) {
  switch (action.type) {
    default:
      return state;
  }
}
