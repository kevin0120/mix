// @flow
import {
  call,
  cancel,
  fork,
  join,
  put,
  race,
  take
} from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import { orderStepUpdateApi } from '../../api/order';
// eslint-disable-next-line import/no-cycle
import { orderActions } from '../order/action';
import { ORDER } from '../order/constants';
import STEP_STATUS from './constants';
// eslint-disable-next-line import/no-cycle
import stepTypes from './stepTypes';
import type { tStepDataReducer, tAnyStepState, tRunSubStepCallbacks } from './interface/typeDef';
import type { tStep } from '../order/interface/typeDef';
import { IWorkStep } from './interface/IWorkStep';

function invalidStepStatus(stepType, status) {
  if (!stepType) {
    throw Error(`invalid stepType ${stepType}`);
  }
  if (!status) {
    throw Error(`trying to invalid status ${status} of ${stepType}`);
  }
  throw Error(`step type ${stepType}  has empty status ${status}`);
}


export default class Step implements IWorkStep {
  _id = '';

  _name = '';

  _desc = '';

  _info = null;

  _type = '';

  _skippable = false;

  _undoable = false;

  _status = STEP_STATUS.READY;

  _stateToRun = STEP_STATUS.ENTERING;

  _payload = {};

  _data = {};

  _times = [];

  _statusTasks = {};

  _runningStatusTask = null;

  _apis = {
    updateStatus: orderStepUpdateApi
  };

  _steps: Array<tClsStep> = [];

  _onLeave = null;

  constructor(stepObj: tStep) {
    this._id = stepObj.id;
    this.update(stepObj);
    this.run = this.run.bind(this);
    this.timerStart = this.timerStart.bind(this);
    this.timerStop = this.timerStop.bind(this);
    this.updateData = this.updateData.bind(this);
  }

  update(stepObj: tStep) {
    this._name = stepObj.name || 'unnamed step';
    this._desc = stepObj.desc || '';
    this._info = stepObj.info;
    this._type = stepObj.type || '';
    this._skippable = stepObj.skippable || false;
    this._undoable = stepObj.undoable || false;
    this._status = stepObj.status || this._status;
    this._steps = stepObj.steps ? stepObj.steps.map(sD => {
      const existStep = this._steps.find(s => s._id === sD.id);
      if (existStep) {
        existStep.update(sD);
        return existStep;
      }
      return new (stepTypes[sD.type](Step))(sD);
    }) : this._steps;

    this._payload = stepObj.payload || this._payload;
  }

  get id() {
    return this._id;
  }

  get name() {
    return this._name;
  }

  get desc() {
    return this._desc;
  }

  get description() {
    return this._desc;
  }

  get type() {
    return this._type;
  }

  get skippable() {
    return this._skippable;
  }

  get undoable() {
    return this._undoable;
  }

  get status() {
    return this._status;
  }

  get steps() {
    return this._steps;
  }

  get payload() {
    return this._payload;
  }

  get data() {
    return this._data;
  }

  timeCost() {
    return ((this._times || []).length % 2 === 0
        ? this._times || []
        : [...this._times, new Date()]
    ).reduce(
      (total, currentTime, idx) =>
        idx % 2 === 0 ? total - currentTime : total - (0 - currentTime),
      0
    );
  }

  timeLost() {
    // TODO: calculate time lost
  }

  timerStart() {
    try {
      const isStarted = this._times.length % 2 === 1;
      if (!isStarted) {
        this._times.push(new Date());
      }
      // yield put(orderActions.updateState());
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step timerStart'
      });
    }
  }

  timerStop() {
    try {
      const isStarted = this._times.length % 2 === 1;
      if (isStarted) {
        this._times.push(new Date());
      }
      // yield put(orderActions.updateState());
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step timerStart'
      });
    }
  }

  * updateData(dataReducer: tStepDataReducer): Saga<void> {
    try {
      this._data = dataReducer(this._data);
      yield put(orderActions.updateState());
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step.UpdateData'
      });
    }
  }

  * _updateStatus({ status, msg }) {
    if (status in this._statusTasks) {
      try {
        this._status = status;
        yield call(this._apis.updateStatus, this.id, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    } else {
      throw new Error(`step ${this._name}(${this._id})has no status ${status}`);
    }
  }

  * _runStatusTask({ status, msg }) {
    if (status in this._statusTasks) {
      try {
        if (this._runningStatusTask) {
          yield cancel(this._runningStatusTask);
        }
        const taskToRun =
          this._statusTasks[status] && this._statusTasks[status].bind(this) ||
          (() => invalidStepStatus(this._type, status));

        this._runningStatusTask = yield fork(
          taskToRun,
          ORDER,
          orderActions,
          msg
        );
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  }

  * run(stateToRun: tAnyStepState = this._stateToRun): Saga<void> {
    const runStatusTask = this._runStatusTask.bind(this);
    const updateStatus = this._updateStatus.bind(this);
    this.timerStart();

    function* runStep() {
      try {
        while (true) {
          const action = yield take(
            a => a.type === ORDER.STEP.STATUS && a.step === this
          );
          yield fork(updateStatus, action);
          yield fork(runStatusTask, action);
        }
      } catch (e) {
        CommonLog.lError(e, {
          at: 'runStep'
        });
      }
    }

    try {
      const step = yield fork([this, runStep]);
      yield put(orderActions.stepStatus(this, stateToRun));
      yield join(step);
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step root'
      });
    } finally {
      this.timerStop();
      if (this._onLeave) {
        yield call([this, this._onLeave]);
      }
      this._runningStatusTask = null;
    }
  }

  * runSubStep(step: tClsStep, callbacks: tRunSubStepCallbacks): Saga<void> {
    try {
      const { exit, next, previous } = yield race({
        exit: call(step.run),
        next: take(
          action =>
            (action.type === ORDER.STEP.FINISH && action.step === step) ||
            action.type === ORDER.STEP.DO_NEXT
        ),
        previous: take(ORDER.STEP.DO_PREVIOUS)
      });
      if (exit && callbacks && callbacks.onExit) {
        yield call(callbacks.onExit);
      }
      if (next && callbacks && callbacks.onNext) {
        yield call(callbacks.onNext);
      }
      if (previous && callbacks && callbacks.onPrevious) {
        yield call(callbacks.onPrevious);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'runSub ',
        id: `${this._id},${step.id}`
      });
    } finally {
      CommonLog.Info(`run substep finished (${this._id}-${this._name})`);
    }
  }
}

export type tClsStep = Step;
