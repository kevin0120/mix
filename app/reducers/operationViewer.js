import { OPERATION_VIEWER } from '../actions/actionTypes';
import { genReducers } from './utils';

const initOperationViewer = {
  list: [],
  detail: {},
  loading:false,
};

const operationViewerReducers={
  [OPERATION_VIEWER.DETAIL_FETCH_OK]: (state, action) => ({
    ...state,
    detail: action.data,
    loading:false
  }),
  [OPERATION_VIEWER.LIST_FETCH_OK]: (state, action) => ({
    list: action.data,
    detail: {},
    loading:false
  }),
  [OPERATION_VIEWER.EDIT_START]: (state, action) => ({
    ...state,
    loading:true
  })
};

export default genReducers(operationViewerReducers,initOperationViewer);
