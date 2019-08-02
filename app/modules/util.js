//@ flow
import type { Saga } from 'redux-saga';
import { fork, take } from 'redux-saga/effects';

export function genReducers(reducers, initState = {}) {
  return (state = initState, action) => {
    if (reducers[action.type]) {
      return reducers[action.type](state, action);
    }
    return state;
  };
}

export function watchWorkers(workers, channel): Saga<void> {
  return function* watcher() {
    try {
      while (true) {
        const action = yield take(channel || Object.keys(workers));
        if (workers[action.type].length > 1) {
          const effect = workers[action.type][0];
          const worker = workers[action.type][1];
          yield effect(worker, action);
        } else {
          yield fork(workers[action.type], action);
        }
      }
    } catch (e) {
      console.error('saga watcher error:', e, workers, channel);
    }
  };
}
