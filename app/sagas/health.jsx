import { call, take } from 'redux-saga/effects';
import { HEALTH } from '../actions/actionTypes';


export function* watchHealth() {
  while (true){
    const { data } = yield take(HEALTH.HEALTH);
    yield call(handleHealth, data)
  }
}

export function* handleHealth(data) {

}
