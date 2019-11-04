// @flow
import type { Saga, Task } from 'redux-saga';
import type { tAnyStepStatus, tRunSubStepCallbacks, tStepDataReducer, tStepPayload, tStep, tStepInfo } from './typeDef';
import type { tCallable } from '../../typeDef';

export interface IWorkStep {
  _code: string,
  _name: string,
  _desc: string,
  _info: ?tStepInfo,
  _type: string,
  _skippable: boolean,
  _undoable: boolean,
  +_onLeave: ?tCallable<Array<any>, any>,
  _payload: ?tStepPayload,
  _statusTasks: { [key: string]: tCallable<Array<any>, void> },
  _runningStatusTask: ?Task<any>,
  _status: tAnyStepStatus,
  _stateToRun: tAnyStepStatus,
  _data: any,
  _times: Array<Date>,
  _steps: Array<IWorkStep>,
  _apis: {
    updateStatus: tCallable<Array<any>, any>
  },
  +constructor: ({ [key: string]: any }, ...Array<any>)=>void,
  +update: (tStep)=>void,
  +code: string,
  +name: string,
  +desc: string,
  +info: ?tStepInfo,
  +type: string,
  +skippable: boolean,
  +undoable: boolean,
  +status: tAnyStepStatus,
  +steps: Array<IWorkStep>,
  +payload: ?tStepPayload,
  +data: any,

  +run: tCallable<Array<any>, any>,

  update(?{ [key: string]: any }): void,

  timeLost(): number,

  timeCost(): number,

  timerStart(): void,

  timerStop(): void,

  updateData(tStepDataReducer): Saga<void>,

  runSubStep(
    step: IWorkStep,
    callbacks: tRunSubStepCallbacks,
    status: tAnyStepStatus
  ): Saga<void>
}
