import { call, put, select } from 'redux-saga/effects';

import {
  OPERATION_VIEWER,
  fetchOperationDetailStart,
  editOperationEnd
} from './action';
import {
  operationDetailApi,
  operationListApi,
  imageEditApi
} from '../../api/operationViewer';
import { watch } from '../indexSaga';
import { setNewNotification } from '../notification/action';

export default watch({
  [OPERATION_VIEWER.LIST_FETCH_START]: fetchOperationList,
  [OPERATION_VIEWER.DETAIL_FETCH_START]: fetchOperationDetail,
  [OPERATION_VIEWER.EDIT_START]: editOperation
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
    yield put(setNewNotification('error', '获取作业信息失败'));
    yield put({
      type: OPERATION_VIEWER.LIST_FETCH_OK,
      data: []
    });
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
    yield put(setNewNotification('error', '获取作业信息失败'));
    yield put({
      type: OPERATION_VIEWER.DETAIL_FETCH_OK,
      data: {}
    });
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
    const { operationID, points, img } = action;
    const state = yield select();
    const { value: odooUrl } = state.setting.page.odooConnection.odooUrl;
    const response = yield call(imageEditApi, odooUrl, operationID, points, img);
    if (response && response.status === 200) {
      yield put(fetchOperationDetailStart(operationID));
      yield put(setNewNotification('info', '作业信息编辑成功'));

    } else {
      yield put(editOperationEnd(false));
      yield put(setNewNotification('error', '作业信息编辑失败'));
    }
  } catch (e) {
    console.error(e);
    yield put(editOperationEnd(false));
    yield put(setNewNotification('error', '作业信息编辑失败'));
  }
}
