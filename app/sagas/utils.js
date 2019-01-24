import { take, call, put, select, fork } from 'redux-saga/effects';

export function genWatcher(workers) {
  return function* watcher(){
    try {
      while (true) {
        const action = yield take(Object.keys(workers));
        yield fork(workers[action.type],action);
      }
    } catch (e) {
      console.log(e);
    }
  }
}
