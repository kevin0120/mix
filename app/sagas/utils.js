import { take, fork } from 'redux-saga/effects';

export function watch(workers, channel) {
  return function* watcher(){
    try {
      while (true) {
        const action = yield take(channel || Object.keys(workers));
        if(workers[action.type].length>1){
          const effect=workers[action.type][0];
          const worker=workers[action.type][1];
          yield effect(worker,action);
        }else{
          yield fork(workers[action.type],action);
        }
      }
    } catch (e) {
      console.log(e);
    }
  }
}
