// @flow
import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';
import type { Saga } from 'redux-saga';

export type StateType = {
  +notify: Object,
  +setting: Object
};

export type tAction<TType, TPayload> = {|
  type: TType,
  ...TPayload
|};

export type tReducer<S, A> = (state: S, action: A) => S;

export type GetState = () => StateType;

// eslint-disable-next-line flowtype/no-weak-types
export type Dispatch = ReduxDispatch<tAction<string, any>>;

// eslint-disable-next-line flowtype/no-weak-types
export type Store = ReduxStore<StateType, tAction<string, any>>;

// eslint-disable-next-line flowtype/no-weak-types
export type tCallable<TP, TR> = (...TP) => TR | Saga<TR>;

// eslint-disable-next-line flowtype/no-weak-types
export type tPredicate<T> = (...Array<T>) => boolean;

// eslint-disable-next-line flowtype/no-weak-types
export type tFnAction<T> = (...Array<T>) => tAction<any, any>;

export type tListener<T> = {
  predicate: tPredicate<T>,
  action: tFnAction<T>
};

/* eslint-disable flowtype/no-weak-types */
export type tListenerObj = {
  listeners: Array<tListener<any>>,
  add: (predicate: tPredicate<any>, action: tFnAction<any>) => tListener<any>,
  remove: (tListener<any>) => Array<tListener<any>>,
  check: (...args: Array<any>) => Array<tAction<any, any>>
};
/* eslint-enable flowtype/no-weak-types */
