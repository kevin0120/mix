import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';

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
export type Store = ReduxStore<GetState, tAction<string, any>>;
