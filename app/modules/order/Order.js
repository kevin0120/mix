import React from 'react';
import { push } from 'connected-react-router';
import { call, put } from 'redux-saga/effects';
import { ORDER_STATUS } from './model';
import { CommonLog, durationString } from '../../common/utils';
import { orderActions } from './action';
import { orderUpdateApi } from '../../api/order';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import STEP_STATUS from '../step/model';

const stepStatus = (status) => {
  switch (status) {
    case STEP_STATUS.FINISHED:
      return '完成';
    case STEP_STATUS.FAIL:
      return '失败';
    default:
      return '未完成';
  }
};

const OrderMixin = (ClsBaseStep) => class ClsOrder extends ClsBaseStep {

  _apis = {
    updateStatus: orderUpdateApi
  };

  _workingIndex = 0;

  _workingID = null;


  constructor(dataObj) {
    super(...arguments);
    this._status = dataObj.status || ORDER_STATUS.TODO;
  }

  get workingStep() {
    return this._steps[this._workingIndex];
  }

  get workingIndex() {
    return this._workingIndex;
  }

  * onNext() {
    try {
      this._workingIndex += 1;
      if (this._workingIndex >= this._steps.length) {
        yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  }

  * onPrevious() {
    if (this._workingIndex - 1 < 0) {
      // yield put(orderActions.finishOrder(this));
    } else {
      this._workingIndex -= 1;
    }
  }

  _statusTasks = {
    * [ORDER_STATUS.TODO]() {

    },
    * [ORDER_STATUS.WIP]() {
      try {
        this._workingIndex = this._workingIndex >= this._steps.length ? 0 : this._workingIndex;
        while (true) {
          CommonLog.Info('Doing Order...', this._workingIndex);
          const step = this.workingStep;
          if (step) {
            yield call([this, this.runSubStep], step, {
              onNext: this.onNext.bind(this),
              onPrevious: this.onPrevious.bind(this)
            });
          } else {
            yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
          }
        }
      } catch (e) {
        CommonLog.lError(e, { at: 'ORDER_STATUS.WIP' });
      } finally {
        CommonLog.Info('order doing finished');
      }
    },
    * [ORDER_STATUS.DONE]() {
      try {
        if (this.workingStep) {
          this.workingStep.timerStop();
        }
        const data = this._steps.map(s => [s.name, durationString(s.timeCost()), stepStatus(s.status)]);
        yield put(
          dialogActions.dialogShow({
            buttons: [
              {
                label: 'Common.Yes',
                color: 'info'
              }
            ],
            closeAction: push('/app'),
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
        const err = (e: Error);
        CommonLog.lError(`showResult error: ${err.message}`);
      } finally {
        CommonLog.Info('order done');
      }
    },
    * [ORDER_STATUS.PENDING]() {
      try {
        this.workingStep.timerStop();
        yield put(orderActions.finishOrder(this));
      } catch (e) {
        CommonLog.lError(e, {
          at: 'ORDER_STATUS.PENDING'
        });
      }
    },
    * [ORDER_STATUS.CANCEL]() {
      try {
        this.workingStep.timerStop();
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

export type tClsOrder = ClsOrder;
