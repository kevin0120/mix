// @flow
import { LOADING } from './action';
import type { tLoadingState } from './typeDef';
import type { tAction } from '../typeDef';

export default function(state: tLoadingState = false, action: tAction<any, any>) {
  switch (action.type) {
    case LOADING.START:
      return true;
    case LOADING.STOP:
      return false;
    default:
      return state;
  }
}
