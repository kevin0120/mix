// @flow
import type { Saga, TakeableChannel, Effect } from 'redux-saga';
import { fork, take } from 'redux-saga/effects';
import type { Action } from 'redux';
import { remove } from 'lodash-es';
import type {
  tAction,
  tListener,
  tFnAction,
  tPredicate,
  tReducer,
  tListenerObj
} from './typeDef';

export function genReducers<TState, TActionTypes>(
  reducers: { [key: TActionTypes]: tReducer<TState, tAction<TActionTypes, any>> },
  initState: TState
) {
  return (state: TState = initState, action: tAction<TActionTypes, any>) => {
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

/* eslint-disable flowtype/no-weak-types */
export function makeListener(
  initListeners: Array<tListener<any>> = []
): tListenerObj {
  const listeners = initListeners;
  
  function addListener(
    predicate: tPredicate<any>,
    action: tFnAction<any>
  ): tListener<any> {
    const listener = {
      predicate,
      action
    };
    listeners.push(listener);
    return listener;
  }
  
  function removeListener(listener: tListener<any>): Array<tListener<any>> {
    const idx = listeners.findIndex(l => l === listener);
    if (idx < 0) {
      return initListeners;
    }
    return remove(listeners, (l: tListener<any>) => l === listener);
  }
  
  function check(...args: Array<any>): Array<tAction<any, any>> {
    const actions: tAction<any, any> = [];
    listeners.forEach(l => {
      if (l.predicate(...args)) {
        actions.push(l.action(...args));
      }
    });
    
    return actions;
  }
  
  return {
    listeners,
    add: addListener,
    remove: removeListener,
    check
  };
}

/* eslint-enable flowtype/no-weak-types */
