import React from 'react';
import Step from '../step/Step';
import { ORDER_STATUS } from './model';
import { CommonLog, durationString } from '../../common/utils';
import { call, put, race, take } from 'redux-saga/effects';
import { ORDER, orderActions } from './action';
import { orderUpdateApi } from '../../api/order';
import dialogActions from '../dialog/action';
import { push } from 'connected-react-router';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import STEP_STATUS from '../step/model';

export default class Order extends Step {
  _apis = {
    updateStatus: orderUpdateApi
  };

  _workingIndex = 0;

  _steps = [];

  _status = ORDER_STATUS.TODO;

  constructor() {
    super(...arguments);
  }

  get workingStep() {
    return this._steps[this._workingIndex];
  }

  get workingIndex() {
    return this._workingIndex;
  }


  _statusTasks = {
    * [ORDER_STATUS.TODO]() {

    },
    * [ORDER_STATUS.WIP]() {
      while (true) {
        CommonLog.Info('Doing Order...');
        const step = this.workingStep;
        if (step) {
          step.timerStart();
          const { next, previous } = yield race({
            exit: call(step.run, STEP_STATUS.ENTERING),
            next: take(ORDER.STEP.DO_NEXT),
            previous: take(ORDER.STEP.DO_PREVIOUS)
          });
          step.timerStop();
          if (next) {
            if (this._workingIndex + 1 >= this._steps.length) {
              yield put(orderActions.finishOrder(this));
            }
            this._workingIndex += 1;
          }
          if (previous) {
            if (this._workingIndex - 1 < 0) {
              yield put(orderActions.finishOrder(this));
            }else{
              this._workingIndex -= 1;
            }
          }
        } else {
          yield put(orderActions.stepStatus(this, ORDER_STATUS.DONE));
        }
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