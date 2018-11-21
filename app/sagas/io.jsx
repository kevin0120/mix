import { call, take, put } from 'redux-saga/effects';
import { IO, OPERATION } from '../actions/actionTypes';
import { IO_FUNCTION } from '../reducers/io';
import { continueOperation } from './operation';

export function* watchIO() {
  while (true) {
    const { data } = yield take(IO.FUNCTION);
    yield call(handleIOFunction, data);
  }
}

export function* handleIOFunction(data) {
  switch (data) {
    case IO_FUNCTION.IN.RESET: {
      // 复位

      yield call(continueOperation);
      break;
    }
    case IO_FUNCTION.IN.BYPASS: {
      // 强制放行

      yield put({ type: OPERATION.FINISHED, data: [] });
      break;
    }
    case IO_FUNCTION.IN.MODE_SELECT: {
      // 模式选择

      break;
    }

    default:
      break;
  }
}
