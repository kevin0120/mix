import { fetchLogo } from './api/logo';
import { select, call, take, put } from 'redux-saga/effects';
import { LOGO } from '../actions/actionTypes';

export default function* logoFlow() {
  try {
    while (true) {
      yield take(LOGO.FETCH_START);
      const state = yield select();
      const { odooUrl } = state.setting.page.odooConnection;
      const resp = yield call(fetchLogo, odooUrl.value);
      if (resp && resp.status === 200 ) {
        const {logo} = resp.data;
        yield put({
          type: LOGO.FETCH_OK,
          logo
        });
      }
    }
  } catch (e) {
    console.log(e);
  }
}
