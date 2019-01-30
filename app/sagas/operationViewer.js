import { call, fork, put, select, take } from 'redux-saga/effects';

import { OPERATION_VIEWER } from '../actions/actionTypes';
import { operationDetailApi, operationListApi } from './api/operationViewer';
import {watch} from './utils';

export default watch({
  [OPERATION_VIEWER.LIST_FETCH_START]:fetchOperationList,
  [OPERATION_VIEWER.DETAIL_FETCH_START]:fetchOperationDetail
})

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
    console.log(e);
  }
}

function* fetchOperationDetail(action) {
  try {
    const {operationID}=action;
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const response = yield call(operationDetailApi, odooUrl.value, operationID);
    yield put({
      type: OPERATION_VIEWER.DETAIL_FETCH_OK,
      data: response.data
    });
  } catch (e) {
    console.log(e);
  }
}
