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

// export default function operationViewer(
//   state: object = initOperationViewer,
//   action: actionType
// ) {
//   switch (action.type) {
//     case OPERATION_VIEWER.DETAIL_FETCH_OK: {
//       return {
//         ...state,
//         detail: action.data
//       };
//     }
//     case OPERATION_VIEWER.LIST_FETCH_OK: {
//       return {
//         list: action.data,
//         detail: {}
//       };
//     }
//     default:
//       return state;
//   }
// }
