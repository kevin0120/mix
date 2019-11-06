// @flow
import { call, cancel, fork, join, put, race, take } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import { orderStepUpdateApi } from '../../api/order';
import { orderActions } from '../order/action';
import { ORDER } from '../order/constants';
import { STEP_STATUS } from './constants';
import stepTypes from './stepTypes';
import type {
  tStepDataReducer,
  tAnyStepStatus,
  tRunSubStepCallbacks
} from './interface/typeDef';
import type { IWorkStep } from './interface/IWorkStep';

function invalidStepStatus(stepType, status) {
  if (!stepType) {
    throw Error(`invalid stepType ${stepType}`);
  }
  if (!status) {
    throw Error(`trying to run invalid status ${status} of ${stepType}`);
  }
  throw Error(`step type ${stepType}  has empty status ${status}`);
}

export default class Step implements IWorkStep {
  _id: number = 0;

  _code = '';

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

  _steps = [];

  _apis = {
    updateStatus: orderStepUpdateApi
  };

  // eslint-disable-next-line flowtype/no-weak-types,no-unused-vars
  constructor(stepObj: { [key: string]: any }, ...rest: Array<any>) {
    this._code = stepObj.code;
    this._id = stepObj.id;
    this.update(stepObj);
    /* eslint-disable flowtype/no-weak-types */
    (this: any).run = this.run.bind(this);
    (this: any).timerStart = this.timerStart.bind(this);
    (this: any).timerStop = this.timerStop.bind(this);
    (this: any).updateData = this.updateData.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  // eslint-disable-next-line flowtype/no-weak-types
  update(stepObj: ?{ [key: string]: any }) {
    const { name, desc, info, skippable, undoable, status1, steps, text, test_type, payload } = stepObj || {};
    this._name = name || text || desc || 'unnamed';
    this._desc = desc || '';
    this._info = info || {};
    // eslint-disable-next-line camelcase
    this._type = test_type || '';
    this._skippable = skippable || false;
    this._undoable = undoable || false;
    this._status = status1 || this._status;
    this._steps = steps
      ? (stepObj &&
      stepObj.steps.map<IWorkStep>(sD => {
        const existStep = this._steps.find(s => s.code === sD.code);
        if (existStep) {
          existStep.update(sD);
          return existStep;
        }
        if (!stepTypes[sD.test_type] || typeof stepTypes[sD.test_type] !== 'function') {
          return new Step(sD);
        }
        return new (stepTypes[sD.test_type](Step))(sD);
      })) ||
      []
      : this._steps;
    if (stepObj && stepObj.data) {
      this._data = (() => {
        try {
          return JSON.parse(stepObj.data) || {};
        } catch (e) {
          CommonLog.lError(e, {
            at: 'step updating data',
            data: stepObj.data,
            code: this._code,
            name: this._name
          });
          return {};
        }
      })();
    }
    this._payload = {
      ...this._payload,
      ...(payload || {})
    };
  }

  get code() {
    return this._code;
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

  get info() {
    return this._info;
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

  timeCost(): number {
    return ((this._times || []).length % 2 === 0
        ? this._times || []
        : [...this._times, new Date()]
    ).reduce(
      (total, currentTime, idx) =>
        idx % 2 === 0 ? total - currentTime : total - (0 - currentTime),
      0
    );
  }

  // eslint-disable-next-line class-methods-use-this
  timeLost() {
    // TODO: calculate time lost
    return 0;
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

  * _updateStatus({ status }) {
    if (status in this._statusTasks) {
      try {
        this._status = status;
        yield call(this._apis.updateStatus, this.id, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    } else {
      throw new Error(`step ${this._name}(${this._code})has no status ${status}`);
    }
  }

  * _runStatusTask({ status, msg }) {
    try {
      if (this._runningStatusTask) {
        yield cancel(this._runningStatusTask);
      }
      const taskToRun =
        (this._statusTasks[status] && this._statusTasks[status].bind(this)) ||
        (() => invalidStepStatus(this._type, status));

      this._runningStatusTask = yield fork(taskToRun, ORDER, orderActions, msg);
    } catch (e) {
      CommonLog.lError(e);
    }
  }

  * run(status: tAnyStepStatus): Saga<void> {
    let statusToRun = status;
    if (!statusToRun) {
      statusToRun = this._status;
    }
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

    // TODO deal with when stateToRun is fail / redo step
    try {
      const step = yield fork([this, runStep]);
      yield put(orderActions.stepStatus(this, statusToRun));
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

  * runSubStep(
    step: IWorkStep,
    callbacks: tRunSubStepCallbacks,
    status: tAnyStepStatus
  ): Saga<void> {
    try {
      const { exit, next, previous } = yield race({
        exit: call([step, step.run], status),
        next: take(
          action => action.type === ORDER.STEP.FINISH && action.step === step
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
        code: `${this._code},${step.code}`
      });
      throw e;
    } finally {
      CommonLog.Info(`run substep finished (${this._code}-${this._name})`);
    }
  }
}
