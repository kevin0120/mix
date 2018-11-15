import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';

export type StateType = {
  +notify: object,
  +setting: object
};

export type Action = {
  +type: string
};

export type GetState = () => StateType;

export type Dispatch = ReduxDispatch<Action>;

export type Store = ReduxStore<GetState, Action>;
