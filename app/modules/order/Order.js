import React from 'react';
import { push } from 'connected-react-router';
import Step from '../step/Step';
import { ORDER_STATUS } from './model';
import { CommonLog, durationString } from '../../common/utils';
import { call, put } from 'redux-saga/effects';
import { orderActions } from './action';
import { orderUpdateApi } from '../../api/order';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';

export default class Order extends Step {
  _apis = {
    updateStatus: orderUpdateApi
  };

  _workingIndex = 0;

  _workingID = null;

  _steps = [];

  _status = ORDER_STATUS.TODO;

  get workingStep() {
    return this._steps[this._workingIndex];
  }

  get workingIndex() {
    return this._workingIndex;
  }

  * onNext() {
    try {
      console.log(this._workingIndex,this._steps.length);
      if (this._workingIndex + 1 >= this._steps.length) {
        yield put(orderActions.finishOrder(this));
      }
      this._workingIndex += 1;
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
        while (true) {
          CommonLog.Info('Doing Order...');
          const step = this.workingStep;
          if (step) {
            yield call(this.runSubStep, step, {
              onNext: this.onNext.bind(this),
              onPrevious: this.onPrevious.bind(this)
            });
          } else {
            yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
          }
        }
      } catch (e) {
        CommonLog.lError(e, { at: ORDER_STATUS.WIP });
      }
    },
    * [ORDER_STATUS.DONE]() {
      try {
        if (this.workingStep) {
          this.workingStep.timerStop();
        }
        const data = this._steps.map(s => [s.name, durationString(s.timeCost())]);
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
                tableHead={['工步名称', '耗时']}
                tableData={data}
                colorsColls={['info']}
              />
            )
          })
        );
        yield put(orderActions.finishOrder());
      } catch (e) {
        const err = (e: Error);
        CommonLog.lError(`showResult error: ${err.message}`);
      }
    },
    * [ORDER_STATUS.PENDING]() {
      this.workingStep.timerStop();
      yield put(orderActions.finishOrder());
    },
    * [ORDER_STATUS.CANCEL]() {
      this.workingStep.timerStop();
      yield put(orderActions.finishOrder());
    }
  };
}