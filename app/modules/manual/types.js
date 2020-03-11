import type { Dispatch as ReduxDispatch, Store as ReduxStore } from 'redux';

export type StateType = {
  result: Array,
  counter: number,
  scanner: string,
  tool: string,
  pset: number,
  working: boolean,
  timeline: object
};
export type counterStateType = {
  +counter: Object
};

export type Action = {
  +type: string,
  +scanner: string,
  +ControllerSN: string,
  +tool: string,
  +pset: number,
  +result: any,
  +story: object

};

export type GetState = () => counterStateType;

export type Dispatch = ReduxDispatch<Action>;

export type Store = ReduxStore<GetState, Action>;
