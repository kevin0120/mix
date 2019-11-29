// @flow
import { call, cancel, fork, join, put, race, take } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import { orderActions } from '../order/action';
import { ORDER } from '../order/constants';
import stepTypes from '../step/stepTypes';
import type {
  tStep,
  tAnyStatus,
  tRunSubStepCallbacks,
  tStepType
} from '../step/interface/typeDef';
import type { IWorkable } from './IWorkable';
import type { tWorkableData } from './typeDef';

export default class Workable implements IWorkable {
  _id = 0;

  get id() {
    return this._id;
  }

  _code = '';

  get code() {
    // if (this._code) {
    return this._code;
    // }
    // throw new Error(`workable(${String(this._code)}) has invalid code`);
  }

  _status = '';

  get status() {
    return this._status;
  }

  _payload = {};

  get payload() {
    return this._payload;
  }

  _data = {};

  get data() {
    return this._data;
  }

  _steps = [];

  get steps() {
    return this._steps;
  }

  _desc = '';

  get desc() {
    return this._desc;
  }

  _statusTasks = {};

  _runningStatusTask = null;

  _times = [];

  // eslint-disable-next-line flowtype/no-weak-types,no-unused-vars
  constructor(workableData: ?tWorkableData) {
    const { code, id } = workableData || {};
    this._code = code || this._code;
    this._id = id || this._id;
    this.update(workableData);
    /* eslint-disable flowtype/no-weak-types */
    (this: any).run = this.run.bind(this);
    (this: any).timerStart = this.timerStart.bind(this);
    (this: any).timerStop = this.timerStop.bind(this);
    (this: any).updateData = this.updateData.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  // eslint-disable-next-line flowtype/no-weak-types
  update(workableData: ?$Shape<tWorkableData>) {
    // console.log(this, workableData);
    const { steps, status, payload, data } = workableData || {};
    this._status = status || this._status;
    this._steps = steps
      ? ((steps: any): Array<tStep>).map<IWorkable>((sD: tStep) => {
      const existStep = this._steps.find(s => s.code === sD.code);
      if (existStep) {
        existStep.update(((sD: any): tWorkableData));
        return existStep;
      }
      const stepType: ?tStepType = (sD: tStep).test_type;
      if (
        !stepType ||
        !stepTypes[stepType] ||
        typeof stepTypes[stepType] !== 'function'
      ) {
        return new Workable(((sD: any): tWorkableData));
      }
      return new (stepTypes[stepType](Workable))(
        ((sD: any): tWorkableData)
      );
    }) || []
      : this._steps;
    if (data) {
      this._data = (() => {
        try {
          return JSON.parse(data) || {};
        } catch (e) {
          CommonLog.lError(e, {
            at: 'step updating data',
            data,
            code: this._code
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

  * updateData(dataReducer: Function): Saga<void> {
    try {
      this._data = dataReducer(this._data);
      yield put(orderActions.updateState());
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step.UpdateData'
      });
    }
  }

  * updateStatus({ status }: { status: tAnyStatus }): Saga<void> {
    const validStatus = Object.keys(this._statusTasks || {});
    if (validStatus.includes(status)) {
      this._status = status;
      try {
        yield put(orderActions.updateState());
      } catch (e) {
        CommonLog.lError(e);
        throw e;
      }
    } else {
      throw new Error(
        `workable (${String(this._code)}) has no status ${status}`
      );
    }
  }

  * _runStatusTask({ status, config }) {
    try {
      if (this._runningStatusTask) {
        yield cancel(this._runningStatusTask);
      }

      if (!this.status) {
        throw new Error(`trying to run invalid status ${status}`);
      }
      const taskToRun = this._statusTasks[status].bind(this);
      this._runningStatusTask = yield fork(taskToRun, ORDER, orderActions, config);
    } catch (e) {
      CommonLog.lError(e);
    }
  }

  * run(status: tAnyStatus, config): Saga<void> {
    let statusToRun = status;
    if (!statusToRun) {
      statusToRun = this._status;
    }
    const runStatusTask = this._runStatusTask.bind(this);
    // const updateStatus = this.updateStatus.bind(this);
    this.timerStart();

    function* runStep() {
      try {
        while (true) {
          const action = yield take(
            a => a.type === ORDER.STEP.STATUS && a.step === this
          );
          // yield fork(updateStatus, action);
          yield fork(runStatusTask, {
            status: action.status,
            config
          });
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
    step: IWorkable,
    callbacks: tRunSubStepCallbacks,
    status: tAnyStatus,
    config
  ): Saga<void> {
    try {
      const { exit, next, previous } = yield race({
        exit: call([step, step.run], status, config),
        next: take(
          action => action.type === ORDER.STEP.FINISH && action.step === step
        ),
        previous: take(ORDER.STEP.DO_PREVIOUS)
      });
      if (exit && callbacks && callbacks.onExit) {
        return yield call(callbacks.onExit);
      }
      if (next && callbacks && callbacks.onNext) {
        return yield call(callbacks.onNext);
      }
      if (previous && callbacks && callbacks.onPrevious) {
        return yield call(callbacks.onPrevious);
      }
    } catch (e) {
      CommonLog.lError(e, {
        at: 'runSub ',
        code: `${String(this?._code)},${step?.code}`
      });
      throw e;
    } finally {
      CommonLog.Info(`run substep finished (${String(this._code)})`);
    }
  }

  // eslint-disable-next-line class-methods-use-this
  _onLeave() {
    CommonLog.Info('workable leaving');
  }
}
