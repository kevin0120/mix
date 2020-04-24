// @flow
import type { Saga } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import type { IWorkStep } from './interface/IWorkStep';
import type { IWorkable } from '../workable/IWorkable';
import { orderStepUpdateApi, stepDataApi } from '../../api/order';
import type { tStep, tStepStatus } from './interface/typeDef';
import { STEP_STATUS, stepTypeKeys } from './constants';
import { stepStatusTasks } from './stepStatusTasks';
import notifierActions from '../Notifier/action';

const StepMixin = (ClsWorkable: Class<IWorkable>) =>
  class Step extends ClsWorkable implements IWorkStep {
    _status = STEP_STATUS.READY;

    _statusTasks = stepStatusTasks;

    constructor(stepData: ?$Shape<tStep>) {
      super(stepData);
      this.update(stepData);
    }

    _failureMsg = '';

    get failureMsg() {
      return this._failureMsg;
    }

    _info = null;

    get info() {
      return this._info;
    }

    _type = '';

    get type() {
      const validTypes = Object.values(stepTypeKeys);
      if (!this._type || !validTypes.includes(this._type)) {
        throw new Error(`Step(${this.code}) has invalid type ${this._type}`);
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

    _consumeProduct = null;

    get consumeProduct() {
      return this._consumeProduct;
    }

    _text = null;

    get text() {
      return this._text;
    }

    _toleranceMax = null;

    get toleranceMax() {
      return this._toleranceMax;
    }

    _toleranceMin = null;

    get toleranceMin() {
      return this._toleranceMin;
    }

    _target = null;

    get target() {
      return this._target;
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
        status,
        consume_product: consumeProduct,
        text,
        tolerance_min,
        tolerance_max,
        target
      } = stepData || {};
      this._sequence = sequence;
      this._type = testType;
      this._image = image;
      this._toleranceMin = tolerance_min;
      this._toleranceMax = tolerance_max;
      this._target = target;
      this._skippable = skippable;
      this._undoable = undoable;
      this._failureMsg = failureMsg;
      this._consumeProduct = consumeProduct;
      this._text = text;
      (this: IWorkable)._status = status || STEP_STATUS.READY;
      (this: IWorkable)._desc = desc;
    }

    * updateStatus({ status }: { status: tStepStatus }): Saga<void> {
      try {
        yield call([this, super.updateStatus], { status });
        yield call(orderStepUpdateApi, this.code, status);
      } catch (e) {
        yield put(notifierActions.enqueueSnackbar(
          'Error', `更新工步状态失败（${e.message}）`, {
            at: 'updateStatus',
            code: this.code,
            status
          }
        ));
        throw(e);
      }
    }

    * _onLeave() {
      try {
        yield call(stepDataApi, this.code, this._data);
      } catch (e) {
        CommonLog.lError(e, {
          at: `step (${String((this: IWorkable)._code)})._onLeave`
        });
      }
    }

    * clearData() {
      // eslint-disable-next-line redux-saga/no-unhandled-errors
      yield call(this.updateData, () => ({}));
    }

    * reset() {
      try {
        yield call([this, this.clearData]);
        yield call(stepDataApi, this.code, this.data);
        this._status = STEP_STATUS.READY;
        yield call(orderStepUpdateApi, this.code, this._status);
      } catch (e) {
        throw e;
      }
    }
  };

export default StepMixin;
