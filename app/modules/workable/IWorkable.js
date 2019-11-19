// @flow
import type { Saga, Task } from 'redux-saga';
import type {
  tAnyStatus,
  tRunSubStepCallbacks,
  tStepPayload
} from '../step/interface/typeDef';
import type { tWorkableData } from './typeDef';

export type tWorkableDataReducer = Function;

export interface IWorkable {
  _id: number;
  +id: number;

  _code: string;
  +code: string;

  _status: tAnyStatus;
  +status: tAnyStatus;

  _payload: ?tStepPayload;
  +payload: ?tStepPayload;

  _data: any;
  +data: any;

  _desc: string;
  +desc: string;

  _steps: Array<IWorkable>;
  +steps: Array<IWorkable>;

  _statusTasks: { [key: string]: (...Array<any>) => any };
  _runningStatusTask: ?Task<any>;

  _times: Array<Date>; // to calculate time lost & time cost

  constructor(?tWorkableData): void;

  update(?tWorkableData): void;

  timeLost(): number;

  timeCost(): number;

  timerStart(): void;

  timerStop(): void;

  _runStatusTask({ status: string, msg: string }): any;

  updateStatus({ status: tAnyStatus }): any;

  run(any): any;

  runSubStep(
    step: IWorkable,
    callbacks: tRunSubStepCallbacks,
    status: tAnyStatus
  ): Saga<void>;

  updateData(tWorkableDataReducer): Saga<void>;

  _onLeave(any): any;
}
