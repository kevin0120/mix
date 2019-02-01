import { TOOLS } from '../actions/actionTypes';
import { genReducers } from './utils';

const initState = {
  toolSN: '',
  status: 'disconnected'
};

const reducers = {
  [TOOLS.STATUS_CHANGE]: updateToolStatus
};

export default genReducers(reducers, initState);

function updateToolStatus(state, action) {
  const { toolSN, status } = action;
  return {
    ...state,
    toolSN,
    status
  };
}
