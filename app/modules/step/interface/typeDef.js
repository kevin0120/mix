// @flow

export type tStepDataReducer = Function;
export type tAnyStepState = string;
export type tRunSubStepCallbacks = {
  onExit?: Function,
  onNext?: Function,
  onPrevious?: Function

};

export type tStepInfo = {
  [key: string]: string | number
};
