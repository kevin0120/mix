import { call, put, select } from 'redux-saga/effects';

import { OPERATION_VIEWER } from '../actions/actionTypes';
import { operationDetailApi, operationListApi, imageEditApi } from './api/operationViewer';
import { watch } from './utils';

export default watch({
  [OPERATION_VIEWER.LIST_FETCH_START]: fetchOperationList,
  [OPERATION_VIEWER.DETAIL_FETCH_START]: fetchOperationDetail
});

function* fetchOperationList() {
  try {
    const state = yield select();
    const { hmiSn, odooUrl } = state.setting.page.odooConnection;
    const response = yield call(operationListApi, odooUrl.value, hmiSn.value);
    yield put({
      type: OPERATION_VIEWER.LIST_FETCH_OK,
      data: response.data
    });
  } catch (e) {
    console.error(e);
  }
}

function* fetchOperationDetail(action) {
  try {
    const { operationID } = action;
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const response = yield call(operationDetailApi, odooUrl.value, operationID);
    yield put({
      type: OPERATION_VIEWER.DETAIL_FETCH_OK,
      data: response.data
    });
  } catch (e) {
    console.error(e);
  }
}

// 图片编辑
// {
//   "points": [
//   {
//     "y_offset": 0,
//     "x_offset": 0,
//     "sequence": 0
//   }
// ],
//   "img": "string"
// }
function* editOperation(action) {
  try {
    const { points, img } = action;
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const response = yield call(imageEditApi, odooUrl, points, img);
    console.log(response);
  } catch (e) {
    console.error(e);
  }
}
