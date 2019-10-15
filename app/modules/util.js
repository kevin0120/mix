// @flow
import type { Saga, TakeableChannel, Effect } from 'redux-saga';
import { fork, take } from 'redux-saga/effects';
import type { Action } from 'redux';
import type { tReducer } from './typeDef';

export function genReducers<TState, TActionTypes>(
  reducers: { [key: TActionTypes]: tReducer<TState, Action<TActionTypes>> },
  initState: TState
) {
  return (state: TState = initState, action: Action<TActionTypes>) => {
    if (reducers[action.type]) {
      return reducers[action.type](state, action);
    }
    return state;
  };
}

export function watchWorkers(
  workers: {
    [key: string]: (any) => Saga<any> | [(any) => Effect, (any) => Saga<any>]
  },
  channel: ?TakeableChannel<any>
) {
  return function* watcher(): Saga<void> {
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
