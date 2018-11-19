import { call, take } from 'redux-saga/effects';
import { IO } from '../actions/actionTypes';
import { IO_FUNCTION } from '../reducers/io';


export function* watchIO() {
  while (true){
    const { data } = yield take(IO.FUNCTION);
    yield call(handleIOFunction, data)
  }
}

export function* handleIOFunction(data) {
  switch(data) {
    case IO_FUNCTION.IN.RESET: {
      break;
    }
    case IO_FUNCTION.IN.BYPASS: {
      break;
    }
    case IO_FUNCTION.IN.MODE_SELECT: {
      break;
    }
  }
}
