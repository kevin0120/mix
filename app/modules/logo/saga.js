import { select, call, put } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import fetchLogo from './api';
import { fetchLogoOK } from './action';
import { CommonLog } from '../../common/utils';

export default function* logoFlow(): Saga<void> {
  try {
    const state = yield select();
    const { odooUrl } = state.setting.page.odooConnection;
    const resp = yield call(fetchLogo, odooUrl.value);
    if (resp?.status === 200) {
      const { logo } = resp.data;
      yield put(fetchLogoOK(logo));
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}
