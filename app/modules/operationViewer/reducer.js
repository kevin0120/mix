import { OPERATION_VIEWER } from './action';
import { genReducers } from '../util';

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
  [OPERATION_VIEWER.EDIT_START]: (state) => ({
    ...state,
    loading:true
  }),
  [OPERATION_VIEWER.EDIT_END]:(state)=>({
    ...state,
    loading:false,
  })
};

export default genReducers(operationViewerReducers,initOperationViewer);
