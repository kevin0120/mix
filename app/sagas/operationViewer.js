import { call, fork, put, select, take } from 'redux-saga/effects';

import { OPERATION_VIEWER } from '../actions/actionTypes';
import { operationDetailApi, operationListApi } from './api/operationViewer';


export default function* watchOperationViewer() {
  try {
    while (true) {
      const action = yield take([
        OPERATION_VIEWER.LIST_FETCH_START,
        OPERATION_VIEWER.DETAIL_FETCH_START
      ]);
      switch (action.type) {
        case OPERATION_VIEWER.LIST_FETCH_START:
          yield fork(fetchOperationList);
          break;
        case OPERATION_VIEWER.DETAIL_FETCH_START:
          yield fork(fetchOperationDetail, action.operationID);
          break;
        default:
          break;
      }
    }
  } catch (e) {
    console.log(e);
  }
}

function* fetchOperationList() {
  try {
    const state = yield select();
    const { hmiSn, odooUrl } = state.setting.page.odooConnection;
    const response = yield call(
      operationListApi,
      odooUrl.value,
      hmiSn.value
    );
    console.log('fetchOperationList:', response);
    yield put({
      type:OPERATION_VIEWER.LIST_FETCH_OK,
      data:response.data
    });
  } catch (e) {
    console.log(e);
  }
}


function* fetchOperationDetail(operationID) {
  try {
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const response = yield call(
      operationDetailApi,
      odooUrl.value,
      operationID
    );
    yield put({
      type:OPERATION_VIEWER.DETAIL_FETCH_OK,
      data:response.data
    });
    console.log('fetchOperationDetail:', response);
  } catch (e) {
    console.log(e);
  }
}

