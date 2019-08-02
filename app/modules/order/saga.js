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
  doable,
  viewingOrder
} from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { durationString, timeCost } from '../../common/utils';
import type { tStep } from './model';

const mapping = {
  onOrderFinish: showResult,
  onOrderView: showOverview
};

function* showResult(order) {
  try {
    const oSteps = orderSteps(order) || [];
    const data = oSteps.map(s => ([
      s.name,
      durationString(timeCost(s.times))
    ]));
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
    console.error('showResult error', e);
  }
}

function* showOverview(order) {
  try {
    const WIPOrder = yield select(s => workingOrder(s.order));
    const vOrderSteps = yield select((state => orderSteps(viewingOrder(state.order))));
    const data = vOrderSteps.map((s: tStep, idx) => ([
      idx + 1,
      s.name,
      s.type,
      s.info
    ]));
    yield put(
      dialogActions.dialogShow({
        buttons: [
          {
            label: 'Common.Close',
            color: 'warning'
          },
          (!WIPOrder) && doable(order) && {
            label: 'Order.Start',
            color: 'info',
            action: orderActions.workOn(order)
          }
        ],
        title: i18n.t('Order.Overview'),
        content: <Table
          tableHeaderColor="info"
          tableHead={[
            i18n.t('Common.Idx'),
            i18n.t('Order.Step.name'),
            i18n.t('Order.Step.type'),
            i18n.t('Order.Step.desc')
          ]}
          tableData={data}
          colorsColls={['info']}
        />
      }));
  } catch (e) {
    console.error('showOverview error', e);
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
    yield call(mapping.onOrderView, order);
  } catch (e) {
    console.error(e);
  }
}

function* workOnOrder() {
  try {
    const wOrder = yield select(s => workingOrder(s.order));
    const { finish } = yield race({
      exit: call(doOrder),
      finish: take(ORDER.FINISH),
      pending: take(ORDER.PENDING),
      cancel: take(ORDER.CANCEL)
    });
    console.log('order finished');
    if (finish) {
      yield call(mapping.onOrderFinish, wOrder);
    }
  } catch (e) {
    console.error(e);
  }
}

function* doOrder() {
  try {
    while (true) {
      console.log('doing order');
      const wOrder = yield select(state => workingOrder(state.order));
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
        const newWOrder = yield select(state => workingOrder(state.order));
        if (workingIndex(newWOrder) >= orderLength(newWOrder)) {
          yield put(orderActions.finishOrder());
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}
