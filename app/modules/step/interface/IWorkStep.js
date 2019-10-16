import type { Saga } from 'redux-saga';
import type { tAnyStepStatus, tRunSubStepCallbacks, tStepDataReducer, tStepPayload, tStep } from './typeDef';

export interface IWorkStep {
  // _id: string,
  // _name: string,
  // _desc: string,
  // _info: { [key: string]: string | number },
  // _type: string,
  // _skippable: boolean,
  // _undoable: boolean,
  // _onLeave: (void)=>void,
  // _payload: tStepPayload,
  // _statusTasks: { [key: string]: tCallable<void> },
  // _runningStatusTask: Task<any>,
  // _status: tAnyStepStatus,
  // _stateToRun: tAnyStepStatus,
  // _data: any,
  // _times: Array<number>,
  // _steps: IWorkStep,
  // _apis: {
  //   updateStatus: tCallable<any>
  // },
  update: (tStep)=>void,
  id: string,
  name: string,
  desc: string,
  info: { [key: string]: string | number },
  type: string,
  skippable: boolean,
  undoable: boolean,
  status: tAnyStepStatus,
  steps: IWorkStep,
  payload: tStepPayload,
  data: any,

  run(tAnyStepStatus): void | Saga<void>,

  timeLost(): number,

  timeCost(): number,

  timerStart(): void,

  timerStop(): void,

  updateData(tStepDataReducer): Saga<void>,

  runSubStep(IWorkStep, tRunSubStepCallbacks): Saga<void>
}
