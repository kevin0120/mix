import { OPERATION_VIEWER } from '../actions/actionTypes';
import { genReducers } from './utils';

const initOperationViewer = {
  list: [],
  detail: {}
};

const operationViewerReducers={
  [OPERATION_VIEWER.DETAIL_FETCH_OK]: (state, action) => ({
    ...state,
    detail: action.data
  }),
  [OPERATION_VIEWER.LIST_FETCH_OK]: (state, action) => ({
    list: action.data,
    detail: {}
  })
};

export default genReducers(operationViewerReducers,initOperationViewer);
