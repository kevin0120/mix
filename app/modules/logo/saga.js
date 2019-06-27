import { select, call, put } from 'redux-saga/effects';
import fetchLogo from '../../api/logo';
import { LOGO } from './action';

export default function* logoFlow() {
  try {
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const resp = yield call(fetchLogo, odooUrl.value);
    if (resp && resp.status === 200) {
      const { logo } = resp.data;
      yield put({
        type: LOGO.FETCH_OK,
        logo
      });
    }
  } catch (e) {
    console.error(e);
  }
}
