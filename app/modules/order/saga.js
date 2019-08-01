// @flow
import { take, call, race, select, put, all, takeLatest, takeEvery } from 'redux-saga/effects';
import { push } from 'connected-react-router';
import React from 'react';
import { ORDER, orderActions } from './action';
import steps from '../step/saga';
import {
  workingStep,
  workingOrder,
  workingIndex,
  stepType,
  orderLength,
  orderSteps,
  doable
} from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { durationString, timeCost } from '../../common/utils';

const mapping = {
  onOrderFinish: showResult
};

function* showResult() {
  try {
    const workingOrderSteps = yield select((state => orderSteps(workingOrder(state.order))));
    const data = workingOrderSteps.map(s => ([
      s.name,
      durationString(timeCost(s.times))
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

export default function* root(): any {
  try {
    yield all([
      takeLatest(ORDER.WORK_ON, workOnOrder),
      takeEvery(ORDER.VIEW, viewOrder)
    ]);
  } catch (e) {
    console.error(e);
  }
}

function* viewOrder({ order }) {
  try {
    const WIPOrder = yield select(s => workingOrder(s.order));
    if ((!WIPOrder) && doable(order)) {
      yield put(orderActions.workOn(order));
    }
  } catch (e) {
    console.error(e);
  }
}

function* workOnOrder() {
  try {
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
  } catch (e) {
    console.error(e);
  }
}

function* doOrder() {
  try {
    while (true) {
      console.log('doing order');
      const wOrder=yield select(state=>workingOrder(state.order));
      const step = workingStep(wOrder);
      const idx = workingIndex(wOrder);
      const type = stepType(step);
      yield put(orderActions.stepTime(idx, new Date()));
      const { next } = yield race({
        exit: call(steps, type),
        next: take(ORDER.STEP.DO_NEXT),
        previous: take(ORDER.STEP.DO_PREVIOUS)
      });
      yield put(orderActions.stepTime(idx, new Date()));
      if (next) {
        const wOrder = yield select(state => workingOrder(state.order));
        if (workingIndex(wOrder) >= orderLength(wOrder)) {
          yield put(orderActions.finishOrder());
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}
