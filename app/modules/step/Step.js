import { call, cancel, fork, join, put, takeEvery } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { orderStepUpdateApi } from '../../api/order';
import { ORDER, orderActions } from '../order/action';


const STEP_STATUS = {
  READY: 'STEP_STATUS_READY',
  ENTERING: 'STEP_STATUS_ENTERING',
  DOING: 'STEP_STATUS_DOING',
  LEAVING: 'STEP_STATUS_LEAVING',
  FAIL: 'STEP_STATUS_FAIL',
  FINISHED: 'STEP_STATUS_FINISHED'
};


function invalidStepStatus(stepType, status) {
  if (!stepType) {
    throw Error(`invalid stepType ${stepType}`);
  }
  if (!status) {
    throw Error(`trying to invalid status ${status} of ${stepType}`);
  }
  throw Error(`step type ${stepType}  has empty status ${status}`);
}

export default class Step {
  _id = '';

  _name = '';

  _desc = '';

  _type = '';

  _skippable = false;

  _undoable = false;

  _status = STEP_STATUS.READY;

  _payload = {};

  _data = {};

  _times = [];

  _statusTasks = {};

  _runningStatusTask = null;

  _apis = {
    updateStatus: orderStepUpdateApi
  };

  _steps = [];

  constructor(stepObj, stepTypes) {
    this._id = stepObj.id;
    this.update(stepObj, stepTypes);

    this.run = this.run.bind(this);
    this.timerStart = this.timerStart.bind(this);
    this.timerStop = this.timerStop.bind(this);
  }

  update(stepObj, stepTypes) {
    this._name = stepObj.name || 'unnamed step';
    this._desc = stepObj.desc || '';
    this._type = stepObj.type || '';
    this._skippable = stepObj.skippable || false;
    this._undoable = stepObj.undoable || false;
    this._status = stepObj.status || this._status;
    this._steps = stepObj.steps?stepObj.steps.map(sD =>{
      const existStep=this._steps.find(s=>s._id===sD.id);
      if(existStep){
        existStep.update(sD,stepTypes);
        return existStep;
      }
      return new stepTypes[sD.type](sD);
    }):this._steps;
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
    return ((this._times || []).length % 2 === 0 ? this._times || [] : [...this._times, new Date()])
      .reduce((total, currentTime, idx) =>
        idx % 2 === 0 ? total - currentTime : total - (0 - currentTime), 0);
  }

  timeLost() {

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

  updateData(dataReducer) {
    this._data = dataReducer(this._data);
  }

  * _updateStatus({ status }) {
    if (status in this._statusTasks) {
      try {
        if (this._runningStatusTask) {
          yield cancel(this._runningStatusTask);
        }
        const taskToRun = this._statusTasks[status]?.bind(this) ||
          (() => invalidStepStatus(this._type, status));

        this._runningStatusTask = yield fork(taskToRun, ORDER, orderActions);

        this._status = status;
        yield call(this._apis.updateStatus, this.id, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    } else {
      throw new Error(`step ${this._name}(${this._id})has no status ${status}`);
    }
  }


  * run(initStatus) {
    const updateStatus = this._updateStatus.bind(this);

    function* runStep() {
      try {
        yield takeEvery((action) =>
          action.type === ORDER.STEP.STATUS && action.step === this,
          updateStatus
        );
      } catch (e) {
        CommonLog.lError(e, {
          at: 'runStep'
        });
      } finally {
        console.log('running task finished');
      }
    }

    try {
      const step = yield fork([this, runStep]);
      yield put(orderActions.stepStatus(this, initStatus));
      yield join(step);
    } catch (e) {
      CommonLog.lError(e, {
        at: 'step root'
      });
    }
  }

}