import { OPERATION_VIEWER } from '../actions/actionTypes';

const initOperationViewer = {
  list: [],
  detail: {}
};
export default function operationViewer(
  state: object = initOperationViewer,
  action: actionType
) {
  switch (action.type) {
    case OPERATION_VIEWER.DETAIL_FETCH_OK: {
      return {
        ...state,
        detail: action.data
      };
    }
    case OPERATION_VIEWER.LIST_FETCH_OK: {
      return {
        list: action.data,
        detail: {}
      };
    }
    default:
      return state;
  }
}
