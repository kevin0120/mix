import type { Saga, Task } from 'redux-saga';
import type { tRunSubStepCallbacks, tStepDataReducer, tStepPayload } from './typeDef';
import type { tStep } from '../../order/interface/typeDef';

export interface IWorkStep {
  _id: string,
  _name: string,
  _desc: string,
  _info: { [key: string]: string | number },
  _type: string,
  _skippable: boolean,
  _undoable: boolean,
  _onLeave: (void)=>void,
  // eslint-disable-next-line flowtype/no-weak-types
  _payload: tStepPayload,
  _statusTasks: { [key: string]: Function },
  _runningStatusTask: Task<any>,
  update: (tStep)=>void,
  id: string,
  name: string,
  desc: string,
  info: { [key: string]: string | number },
  type: string,
  skippable: boolean,
  undoable: boolean,
  run: (void)=> void | Saga<void>,
  timeLost: (void)=>number,
  timeCost: (void)=>number,
  timerStart: (void)=>void,
  timerStop: (void)=>void,
  updateData: (tStepDataReducer)=>Saga<void>,
  runSubStep: (IWorkStep, tRunSubStepCallbacks)=> Saga<void>
}
