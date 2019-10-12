
import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';

export type StateType = {
  +notify: Object,
  +setting: Object
};

export type Action = {
  +type: string
};

export type GetState = () => StateType;

export type Dispatch = ReduxDispatch<Action>;

export type Store = ReduxStore<GetState, Action>;
