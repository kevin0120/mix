// @flow
import type { Saga } from 'redux-saga';
import React from 'react';
import { push } from 'connected-react-router';
import { call, put, select } from 'redux-saga/effects';
import { some, filter } from 'lodash-es';
import { ORDER_STATUS } from './constants';
import { CommonLog, durationString } from '../../common/utils';
import { orderActions } from './action';
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

const stepStatus = status => {
  switch (status) {
    case STEP_STATUS.FINISHED:
      return '完成';
    case STEP_STATUS.FAIL:
      return '失败';
    default:
      return '未完成';
  }
};

const OrderMixin = (ClsBaseStep: Class<IWorkable>) =>
  class ClsOrder extends ClsBaseStep implements IOrder {
    _workingIndex = 0;

    get workingIndex() {
      return this._workingIndex;
    }

    _status = ORDER_STATUS.TODO;

    get status() {
      return this._status;
    }

    _trackCode = '';

    get trackCode() {
      return this._trackCode;
    }

    _productCode = '';

    get productCode() {
      return this._productCode;
    }

    get failSteps(): Array<IWorkStep> {
      const ret = filter(
        this.steps,
        (step: IWorkStep) => step.status === STEP_STATUS.FAIL
      );
      return ((ret: any): Array<IWorkStep>);
    }

    hasFailWorkStep(): boolean {
      return some(
        this.steps,
        (step: IWorkable) => step.status === STEP_STATUS.FAIL
      );
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

    // eslint-disable-next-line flowtype/no-weak-types
    constructor(dataObj: ?$Shape<tOrder>): void {
      // eslint-disable-next-line prefer-rest-params
      super(dataObj);
      this.update.call(this, dataObj);
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

    get workingStep() {
      return (((this: IWorkable)._steps[this._workingIndex]: any): IWorkStep);
    }

    *updateStatus({ status }: { status: tOrderStatus }): Saga<void> {
      try {
        yield call([this, super.updateStatus], { status });
        yield call(orderUpdateApi, this.id, status);
      } catch (e) {
        CommonLog.lError(e);
      }
    }

    _statusTasks = {
      *[ORDER_STATUS.TODO]() {
        try {
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
      *[ORDER_STATUS.WIP]() {
        try {
          this._workingIndex =
            this._workingIndex >= this._steps.length ? 0 : this._workingIndex;
          let status = null;

          const _onPrevious = () => {
            if (this._workingIndex - 1 >= 0) {
              this._workingIndex -= 1;
            }
            status = STEP_STATUS.READY;
          };

          const _onNext = () => {
            this._workingIndex += 1;
          };

          while (true) {
            CommonLog.Info(
              `Doing Order (${this.code}),at ${this.workingIndex} step (${this.workingStep?.code}) `
            );
            const step = this.workingStep;
            if (step) {
              console.log(status);
              yield call(
                [this, this.runSubStep],
                step,
                {
                  onNext: _onNext,
                  onPrevious: _onPrevious
                },
                status
              );
            } else {
              yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
            }
            status = null;
          }
        } catch (e) {
          CommonLog.lError(e, { at: 'ORDER_STATUS.WIP' });
          yield put(notifyActions.enqueueSnackbar('Error', e.message));
          yield put(orderActions.stepStatus(this, ORDER_STATUS.PENDING));
        } finally {
          CommonLog.Info('order doing finished');
        }
      },
      *[ORDER_STATUS.DONE]() {
        try {
          if (this._workingIndex > this._steps.length - 1) {
            this._workingIndex = this._steps.length - 1;
          }
          if (this._workingIndex < 0) {
            this._workingIndex = 0;
          }
          const data = this._steps.map(s => [
            s.code,
            durationString(s.timeCost()),
            stepStatus(s.status)
          ]);
          const { reportFinish } = yield select(s => s.setting.systemSettings);
          let confirm = {
            label: 'Common.Yes',
            color: 'info'
          };
          let closeAction = push('/app');
          if (reportFinish) {
            const code = this._id;
            const trackCode = '';
            const workCenterCode = yield select(s => s.systemInfo.workcenter);
            const productCode = '';
            const dateComplete = new Date();
            const { operation } = this.payload;
            closeAction = [
              orderActions.reportFinish(
                code,
                trackCode,
                productCode,
                workCenterCode,
                dateComplete,
                operation
              ),
              push('/app')
            ];
            confirm = {
              label: '完工',
              color: 'info'
            };
          }
          yield put(
            dialogActions.dialogShow({
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
      *[ORDER_STATUS.PENDING]() {
        try {
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.PENDING'
          });
        }
      },
      *[ORDER_STATUS.CANCEL]() {
        try {
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.CANCEL'
          });
        }
      }
    };
  };

export default OrderMixin;
