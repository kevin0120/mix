// @flow
import type { Saga } from 'redux-saga';
import { call } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import type { IWorkStep } from './interface/IWorkStep';
import type { IWorkable } from '../workable/IWorkable';
import { orderStepUpdateApi } from '../../api/order';
import type { tStep, tStepStatus } from './interface/typeDef';
import { STEP_STATUS, stepTypeKeys } from './constants';

const StepMixin = (ClsWorkable: Class<IWorkable>) =>
  class Step extends ClsWorkable implements IWorkStep {
    _failureMsg = '';

    get failureMsg() {
      return this._failureMsg;
    }

    _status = STEP_STATUS.READY;

    _info = null;

    get info() {
      return this._info;
    }

    _type = '';

    get type() {
      const validTypes = Object.values(stepTypeKeys);
      if (!this._type || !validTypes.includes(this._type)) {
        throw new Error(`Step(${this.code}) has invalid type${this._type}`);
      }
      return this._type;
    }

    _image = '';

    get image() {
      return this._image;
    }

    _skippable = false;

    get skippable() {
      return this._skippable;
    }

    _undoable = false;

    get undoable() {
      return this._undoable;
    }

    _steps = [];

    get steps() {
      return this._steps;
    }

    _sequence = 0;

    get sequence() {
      return this._sequence;
    }

    constructor(stepData: ?$Shape<tStep>) {
      super(stepData);
      this.update(stepData);
    }

    // eslint-disable-next-line class-methods-use-this
    update(stepData: ?$Shape<tStep>) {
      if (!stepData) {
        return;
      }
      super.update(stepData);
      const {
        sequence,
        test_type: testType,
        failure_msg: failureMsg,
        desc,
        image,
        skippable,
        undoable,
        data,
        status
      } = stepData || {};
      this._sequence = sequence;
      this._type = testType;
      this._image = image;
      this._skippable = skippable;
      this._undoable = undoable;
      this._failureMsg = failureMsg;
      (this: IWorkable)._status = status || STEP_STATUS.READY;
      (this: IWorkable)._data = data;
      (this: IWorkable)._desc = desc;
    }

    *updateStatus({ status }: { status: tStepStatus }): Saga<void> {
      try {
        yield call([this, super.updateStatus], { status });
        yield call(orderStepUpdateApi, this.id, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  };

export default StepMixin;
