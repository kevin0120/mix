// @flow
import type { Saga, Task } from 'redux-saga';
import type { tAnyStepStatus, tRunSubStepCallbacks, tStepDataReducer, tStepPayload, tStep, tStepInfo } from './typeDef';
import type { tCallable } from '../../typeDef';
import type { tOrder } from '../../order/interface/typeDef';

export interface IWorkStep {
  _id: number,
  _name: string,
  _desc: string,
  _info: ?tStepInfo,
  _type: string,
  _skippable: boolean,
  _undoable: boolean,
  +_onLeave: ?tCallable<any>,
  _payload: ?tStepPayload,
  _statusTasks: { [key: string]: tCallable<void> },
  _runningStatusTask: ?Task<any>,
  _status: tAnyStepStatus,
  _stateToRun: tAnyStepStatus,
  _data: any,
  _times: Array<Date>,
  _steps: Array<IWorkStep>,
  _apis: {
    updateStatus: tCallable<any>
  },
  +constructor: ({[key: string]: any}, ...Array<any>)=>void,
  +update: (tStep)=>void,
  +id: number,
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

  +run: tCallable<any>,

  update( {[key:string]: any}): void,

  timeLost(): number,

  timeCost(): number,

  timerStart(): void,

  timerStop(): void,

  updateData(tStepDataReducer): Saga<void>,

  runSubStep(IWorkStep, tRunSubStepCallbacks): Saga<void>
}
