// @flow
import React from 'react';
import { push } from 'connected-react-router';
import { call, put, select } from 'redux-saga/effects';
import { ORDER_STATUS } from './constants';
import { CommonLog, durationString } from '../../common/utils';
import { orderActions } from './action';
import { orderReportStartApi, orderUpdateApi } from '../../api/order';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { STEP_STATUS } from '../step/constants';
import { IOrder } from './interface/IOrder';
import type { IWorkStep } from '../step/interface/IWorkStep';
import loadingActions from '../loading/action';
import notifyActions from '../Notifier/action';

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

const OrderMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsOrder extends ClsBaseStep implements IOrder {
    _apis = {
      updateStatus: orderUpdateApi
    };

    _workingIndex = 0;

    _stateToRun = ORDER_STATUS.TODO;

    _workingID = null;

    _status = ORDER_STATUS.TODO;

    _plannedDateTime = null;

    _trackCode = '';

    _productCode = '';

    get productCode() {
      return this._productCode;
    }

    _workcenter: string = '';

    get workcenter() {
      return this._workcenter;
    }

    _datePlannedStart: Date = new Date();

    get datePlannedStart() {
      return this._datePlannedStart;
    }

    _datePlannedComplete: Date = new Date();

    get datePlannedComplete() {
      return this._datePlannedComplete;
    }

    _productTypeImage: string = '';

    get productTypeImage() {
      return this._productTypeImage;
    }


    // eslint-disable-next-line flowtype/no-weak-types
    constructor(dataObj: { [key: string]: any }, ...rest: Array<any>) {
      super(dataObj, ...rest);
      this._status = dataObj.status || this._status;
      this._trackCode = dataObj.track_code;
      this._productCode = dataObj.product_code;
      this._workcenter = dataObj.workcenter;
      this._datePlannedStart = new Date(dataObj.date_planned_start);
      this._datePlannedComplete = new Date(dataObj.datePlannedComplete);
      this._productTypeImage = dataObj.product_type_image;
    }

    get plannedDateTime() {
      return this._plannedDateTime;
    }

    get workingStep() {
      return (this: IWorkStep)._steps[this._workingIndex];
    }

    get workingIndex() {
      return this._workingIndex;
    }

    _statusTasks = {
      * [ORDER_STATUS.TODO]() {
        try {
          const { reportStart } = yield select(s => s.setting.systemSettings);
          // TODO 开工自检

          if (reportStart) {
            yield put(loadingActions.start());
            const orderCode = this._id;
            const trackCode = this._trackCode;
            const productCode = this._productCode;
            const dateStart = new Date();
            const workCenterCode = yield select(
              s => s.systemInfo.workcenter
            );
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
      * [ORDER_STATUS.WIP]() {
        try {
          this._workingIndex =
            this._workingIndex >= this._steps.length ? 0 : this._workingIndex;
          let status = null;

          const _onPrevious = () => {
            // if (this.workingStep.status === STEP_STATUS.DOING) {
            //
            // }
            if (this._workingIndex - 1 < 0) {
              // yield put(orderActions.finishOrder(this));
            } else {
              this._workingIndex -= 1;
              status = STEP_STATUS.READY;
            }
          };

          const _onNext = () => {
            this._workingIndex += 1;
          };

          while (true) {
            CommonLog.Info('Doing Order...', this._workingIndex);
            const step = this.workingStep;
            if (step) {
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
          const data = this._steps.map(s => [
            s.name,
            durationString(s.timeCost()),
            stepStatus(s.status)
          ]);
          const { reportFinish } = yield select(s => s.setting.systemSettings);
          let confirm = {
            label: 'Common.OK',
            color: 'info'
          };
          let closeAction = push('/app');
          if (reportFinish) {
            const code = this._id;
            const trackCode = '';
            const workCenterCode = yield select(
              s => s.systemInfo.workcenter
            );
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
      * [ORDER_STATUS.PENDING]() {
        try {
          yield put(orderActions.finishOrder(this));
        } catch (e) {
          CommonLog.lError(e, {
            at: 'ORDER_STATUS.PENDING'
          });
        }
      },
      * [ORDER_STATUS.CANCEL]() {
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
