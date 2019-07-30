import { take, call, race, fork, select, put, join, all } from 'redux-saga/effects';
import { push } from 'connected-react-router';
import React from 'react';
import moment from 'moment';
import { ORDER, orderActions } from './action';
import steps from '../step/saga';
import {
  processingStep,
  processingIndex,
  stepType,
  orderLength, orderSteps
} from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { durationString } from '../../common/utils';

const mapping = {
  onOrderFinish: showResult,
};

function* showResult() {
  try {
    const currentOrderSteps = yield select((state => orderSteps(state.order)));
    const data = currentOrderSteps.map(s => ([
      s.name,
      durationString((s.endTime && s.startTime) ? s.endTime - s.startTime : 0)
    ]));
    yield put(
      dialogActions.showDialog({
        hasOk: true,
        closeAction: push('/app'),
        title: i18n.t('Common.Result'),
        content: <Table
          tableHeaderColor="info"
          tableHead={[
            '工步名称',
            '耗时'
          ]}
          tableData={data}
          colorsColls={['info']}
        />
      })
    );
  } catch (e) {
    console.error('returnHome error', e);
  }
}

export default function* root() {
  try {
    while (true) {
      const { order } = yield take(ORDER.TRIGGER);
      const didSwitchOrder = yield fork(function* orderDidSwitch() {
        yield take(ORDER.SWITCH);
      });
      yield put(orderActions.switchOrder(order));
      yield join(didSwitchOrder);
      const { finish } = yield race({
        exit: call(doOrder),
        finish: take(ORDER.FINISH),
        pending: take(ORDER.PENDING),
        cancel: take(ORDER.CANCEL)
      });
      console.log('order finished');
      if (finish) {
        yield call(mapping.onOrderFinish);
      }
    }
  } catch (e) {
    console.error(e);
  }
}

function* doOrder() {
  try {
    while (true) {
      console.log('doing order');
      const step = yield select(state => processingStep(state.order));
      const idx = yield select(state => processingIndex(state.order));
      const type = stepType(step);
      yield put(orderActions.stepStartTime(idx, new Date()));
      yield put(orderActions.stepEndTime(idx, null));

      const { next } = yield race({
        exit: call(steps, type),
        next: take(ORDER.STEP.DO_NEXT),
        previous: take(ORDER.STEP.DO_PREVIOUS)
      });
      yield put(orderActions.stepEndTime(idx, new Date()));
      if (next) {
        const order = yield select(state => state.order);
        if (processingIndex(order) >= orderLength(order)) {
          yield put(orderActions.finishOrder());
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}
