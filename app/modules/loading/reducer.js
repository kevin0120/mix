import { LOADING } from './action';

export default function(state = false, action) {
  switch (action.type) {
    case LOADING.START:
      return true;
    case LOADING.STOP:
      return false;
    default:
      return state;
  }
}
