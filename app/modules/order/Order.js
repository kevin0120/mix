// @flow
import type { Saga } from 'redux-saga';
import React from 'react';
import { push } from 'connected-react-router';
import { call, put, select } from 'redux-saga/effects';
import { filter, some } from 'lodash-es';
import { Typography } from '@material-ui/core';
import { ORDER_STATUS } from './constants';
import { CommonLog, durationString } from '../../common/utils';
import { orderReportStartApi, orderUpdateApi } from '../../api/order';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { STEP_STATUS } from '../step/constants';
import { IOrder } from './interface/IOrder';
import loadingActions from '../loading/action';
import notifyActions from '../Notifier/action';
import type { tOrder, tOrderStatus } from './interface/typeDef';
import type { IWorkable } from '../workable/IWorkable';
import type { IWorkStep } from '../step/interface/IWorkStep';
import { workModes } from '../workCenterMode/constants';
import { ioOutputGroups } from '../io/constants';
import { orderActions } from './action';
import io from '../io';

const stepStatus = status => {
  switch (status) {
    case STEP_STATUS.FINISHED:
      return <Typography color="primary">完成</Typography>;
    case STEP_STATUS.FAIL:
      return <Typography color="error">失败</Typography>;
    default:
      return <Typography>未完成</Typography>;
  }
};

function* redoOrder(step, point) {
  try {
    let redoStep = step;
    if (!redoStep) {
      redoStep = this.steps.find(s => s.status === STEP_STATUS.FAIL);
    }
    if (!redoStep) {
      yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
    }
    // while (true) {
    yield call(
      [this, this.runSubStep],
      redoStep,
      {
        onNext: () => {
        },
        onPrevious: () => {
        }
      },
      STEP_STATUS.READY,
      {
        reworkConfig: {
          point
        }
      }
    );
    if (this._steps.some(s => s.status === STEP_STATUS.FAIL)) {
      yield put(orderActions.stepStatus(this, ORDER_STATUS.FAIL));
    } else {
      yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
    }    // }
  } catch (e) {
    CommonLog.lError(e);
  }
}

const OrderMixin = (ClsBaseStep: Class<IWorkable>) =>
  class ClsOrder extends ClsBaseStep implements IOrder {
    // eslint-disable-next-line flowtype/no-weak-types
    constructor(dataObj: ?$Shape<tOrder>): void {
      // eslint-disable-next-line prefer-rest-params
      super(dataObj);
      this.update.call(this, dataObj);
    }

    _workingIndex = 0;

    get workingIndex() {
      if (this._workingIndex) {
        return this._workingIndex;
      }
      const idx = this._steps.findIndex(s => s.status === STEP_STATUS.DOING);
      return idx >= 0 ? idx : 0;
    }

    set workingIndex(val) {
      this._workingIndex = val;
    }

    _status = ORDER_STATUS.TODO;

    get status() {
      return this._status;
    }

    _data = {};

    get data() {
      return this._data;
    }

    _trackCode = '';

    get trackCode() {
      return this._trackCode;
    }

    _productCode = '';

    get productCode() {
      return this._productCode;
    }

    _statusTasks = {
      * [ORDER_STATUS.TODO](config) {
        try {
          yield put(io.action.setIOOutput({ group: ioOutputGroups.doing, status: true }));
          const { workCenterMode } = yield select();
          if (workCenterMode === workModes.reworkWorkCenterMode) {
            yield put(orderActions.stepStatus(this, ORDER_STATUS.WIP, config));
            return;
          }
          const { reportStart } = yield select(s => s.setting.systemSettings);
          // TODO 开工自检

          if (reportStart) {
            yield put(loadingActions.start());
            const orderCode = this.code;
            const trackCode = this._trackCode;
            const productCode = this._productCode;
            const dateStart = new Date();
            const workCenterCode = yield select(s => s.systemInfo.workcenter);
            const { resources } = this._payload.operation || {};
            // eslint-disable-next-line flowtype/no-weak-types
            yield call(
              orderReportStartApi,
              orderCode,
              trackCode,
              workCenterCode,
              productCode,
              dateStart,
              resources
            );
            yield put(loadingActions.stop());
          }
          yield put(orderActions.stepStatus(this, ORDER_STATUS.WIP));
        } catch (e) {
          yield put(loadingActions.stop());
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.TODO',
            code: this._code,
            name: this._name
          });
          yield put(notifyActions.enqueueSnackbar('Error', e.message));
          yield put(orderActions.stepStatus(this, ORDER_STATUS.PENDING));
        }
      },
      * [ORDER_STATUS.WIP](config = {}) {
        try {
          this.workingIndex =
            this.workingIndex >= this._steps.length ? 0 : this.workingIndex;
          const { step } = config;
          if (step) {
            const stepIndex = this.steps.findIndex(s => s.code === step.code);
            if (stepIndex >= 0) {
              this.workingIndex = stepIndex;
            }
          }

          const mode = yield select(s => s.workCenterMode);
          if (mode === workModes.reworkWorkCenterMode) {
            const reworkConfig = config?.reworkConfig || {};
            const { point } = reworkConfig;
            yield call([this, redoOrder], step, point);
            return;
          }
          let status = null;
          while (true) {
            CommonLog.Info(
              `Doing Order (${this.code}),at ${this.workingIndex} step (${this.workingStep?.code}) `
            );
            const wStep = this.workingStep;
            if (wStep) {
              status = yield call(
                [this, this.runSubStep],
                wStep,
                {
                  onNext: this._onNextStep.bind(this),
                  onPrevious: this._onPreviousStep.bind(this)
                },
                window.debugSettings && window.debugSettings.disableOrderTriggerLimit ?
                  STEP_STATUS.READY :
                  status || (wStep.status === STEP_STATUS.FINISHED ? null : STEP_STATUS.READY)
              );
            } else if (this._steps.some(s => s.status === STEP_STATUS.FAIL)) {
              yield put(orderActions.stepStatus(this, ORDER_STATUS.FAIL));
            } else {
              yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
            }
          }
        } catch (e) {
          CommonLog.lError(e, { at: 'ORDER_STATUS.WIP' });
          yield put(notifyActions.enqueueSnackbar('Error', e.message));
          yield put(orderActions.stepStatus(this, ORDER_STATUS.PENDING));
        } finally {
          CommonLog.Info('order doing finished');
        }
      },
      * [ORDER_STATUS.DONE]() {
        try {
          if (this.workingIndex > this._steps.length - 1) {
            this.workingIndex = this._steps.length - 1;
          }
          if (this.workingIndex < 0) {
            this.workingIndex = 0;
          }
          yield put(io.action.setIOOutput({ group: ioOutputGroups.ready, status: true }));
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.DONE'
          });
          // yield put(orderActions.stepStatus(this, ORDER_STATUS.PENDING));
        } finally {
          CommonLog.Info('order done');
        }
      },
      * [ORDER_STATUS.FAIL]() {
        try {
          yield put(io.action.setIOOutput({ group: ioOutputGroups.error, status: true }));
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.PENDING'
          });
        }
      },
      * [ORDER_STATUS.PENDING]() {
        try {
          yield put(io.action.setIOOutput({ group: ioOutputGroups.warning, status: true }));
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.PENDING'
          });
        }
      },
      * [ORDER_STATUS.CANCEL]() {
        try {
          yield put(io.action.setIOOutput({ group: ioOutputGroups.warning, status: true }));
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.CANCEL'
          });
        }
      }
    };

    _components = [];

    get components() {
      return this._components;
    }

    _workcenter: string = '';

    get workcenter() {
      return this._workcenter;
    }

    _datePlannedStart = null;

    get datePlannedStart() {
      return this._datePlannedStart;
    }

    _datePlannedComplete = null;

    get datePlannedComplete() {
      return this._datePlannedComplete;
    }

    _productTypeImage: string = '';

    get productTypeImage() {
      return this._productTypeImage;
    }

    get workingStep() {
      return (((this: IWorkable)._steps[this.workingIndex]: any): IWorkStep);
    }

    get failSteps(): Array<IWorkStep> {
      const ret = filter(
        this.steps,
        (step: IWorkStep) => step.status === STEP_STATUS.FAIL
      );
      return ((ret: any): Array<IWorkStep>);
    }

    * _onPreviousStep() {
      try {
        const wStep = this.workingStep;
        yield call([wStep, wStep.clearData]);
        const mode = yield select(s => s.workCenterMode);
        if (this.workingIndex - 1 >= 0) {
          if (mode === workModes.normWorkCenterMode) {
            this.workingIndex -= 1;
            yield put(orderActions.jumpToStep(this.workingIndex));
            const nextStep = this.workingStep;
            yield call([nextStep, nextStep.clearData]);
          }
          return STEP_STATUS.READY;
        }
        return null;
      } catch (e) {
        CommonLog.lError(e);
      }
    }

    * _onNextStep() {
      try {
        const mode = yield select(s => s.workCenterMode);
        if (mode === workModes.normWorkCenterMode) {
          this.workingIndex = this.workingIndex + 1;
          yield put(orderActions.jumpToStep(this.workingIndex));
        }
      } catch (e) {
        CommonLog.lError(e, { at: 'order._onNextStep' });
        throw(e);
      }

    }

    clearData() {
      this._steps = [];
      this.workingIndex = 0;
    }

    * _onLeave() {
      try {
        const { workCenterMode } = yield select();
        switch (workCenterMode) {
          case workModes.normWorkCenterMode: {
            const { reportFinish } = yield select(s => s.setting.systemSettings);

            let closeAction = push('/app');
            let confirm = {
              label: 'Common.Yes',
              color: 'info'
            };
            const data = (this: IWorkable)._steps.map(s => [
              s.code,
              durationString(s.timeCost()),
              stepStatus(s.status)
            ]);
            if (reportFinish && (this.status === ORDER_STATUS.FAIL || this.status === ORDER_STATUS.DONE)) {
              closeAction = [
                orderActions.reportFinish(this),
                push('/app')
              ];
              confirm = {
                label: '完工',
                color: 'info'
              };
            }
            yield put(
              dialogActions.dialogShow({
                maxWidth: 'md',
                buttons: [confirm],
                closeAction,
                title: i18n.t('Common.Result'),
                content: (
                  <Table
                    tableHeaderColor="info"
                    tableHead={['工步名称', '耗时', '结果']}
                    tableData={data}
                    colorsColls={['info']}
                  />
                )
              })
            );
            break;
          }
          case workModes.reworkWorkCenterMode: {
            break;
          }
          default:
            break;
        }
      } catch (e) {
        CommonLog.lError(e, { at: 'order onLeave', code: (this: IWorkable)._code });
      }
    }

    update(dataObj: ?$Shape<tOrder>) {
      super.update.call(this, dataObj);
      const {
        status,
        track_code: trackCode,
        product_code: productCode,
        workcenter,
        date_planned_start: datePlannedStart,
        date_planned_complete: datePlannedComplete,
        product_type_image: productTypeImage,
        payload
      } = dataObj || {};

      this._status = status || ORDER_STATUS.TODO;
      this._trackCode = trackCode || '';
      this._productCode = productCode || '';
      this._workcenter = workcenter || '';
      this._datePlannedStart = datePlannedStart
        ? new Date(datePlannedStart)
        : null;
      this._datePlannedComplete = datePlannedComplete
        ? new Date(datePlannedComplete)
        : null;
      this._productTypeImage = productTypeImage || '';
      (this: IWorkable)._desc = payload?.operation?.desc || '';
    }

    * updateStatus({ status }: { status: tOrderStatus }): Saga<void> {
      try {
        yield call([this, super.updateStatus], { status });
        yield call(orderUpdateApi, this.code, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    }

    hasFailWorkStep(): boolean {
      return some(
        this.steps,
        (step: IWorkable) => step.status === STEP_STATUS.FAIL
      );
    }

    hasUnsuccessWorkStep(): boolean {
      return this.steps.some(
        (step: IWorkable) => step.status !== STEP_STATUS.FINISHED
      );
    }
  };

export default OrderMixin;
