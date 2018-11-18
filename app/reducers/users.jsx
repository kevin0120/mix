import { USER } from '../actions/actionTypes';

const defaultUsers = {
  uuid: '',
  name: '',
  avatar: '',
};

type actionType = {
  +type: string
};

export default function users(state: object = defaultUsers, action: actionType) {
  switch (action.type) {
    case USER.FETCH_OK:
      return NewUser(state, action.data);
    default:
      return state;
  }
}

export function NewUser(state, data) {
  return {
    ...state,
    uuid: data.uuid,
    name: data.name,
    avatar: data.avatar,
  };
}
