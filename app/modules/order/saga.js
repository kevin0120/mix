// @flow
import {
  take,
  call,
  race,
  select,
  put,
  all,
  takeLatest,
  takeEvery,
  takeLeading,
  delay
} from 'redux-saga/effects';
import { push } from 'connected-react-router';
import React from 'react';
import type { Saga } from 'redux-saga';
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
import { CommonLog, durationString, timeCost } from '../../common/utils';
import type { tOrder, tStep, tStepArray } from './model';
import type { tCommonActionType } from '../../common/type';

const mapping = {
  onOrderFinish: showResult,
  onOrderView: showOverview
};

function* showResult(order) {
  try {
    const oSteps = orderSteps(order) || [];
    const data = oSteps.map(s => [s.name, durationString(timeCost(s.times))]);
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
  } catch (e) {
    const err = (e: Error);
    CommonLog.lError(`showResult error: ${err.message}`);
  }
}

function* showOverview(order: tOrder) {
  try {
    const WIPOrder: tOrder = yield select(s => workingOrder(s.order));
    const vOrderSteps: ?tStepArray = yield select(state =>
      orderSteps(viewingOrder(state.order))
    );
    const data =
      (vOrderSteps &&
        vOrderSteps.map((s: tStep, idx) => [
          idx + 1,
          s.name,
          s.type,
          s.info
        ])) ||
      [];
    if (WIPOrder === order) {
      // 进行中的工单不显示概览对话框
      return;
    }
    yield put(
      dialogActions.dialogShow({
        buttons: [
          {
            label: 'Common.Close',
            color: 'warning'
          },
          !WIPOrder &&
            doable(order) && {
              label: 'Order.Start',
              color: 'info',
              action: orderActions.workOn(order)
            }
        ],
        title: i18n.t('Order.Overview'),
        content: (
          <Table
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
        )
      })
    );
  } catch (e) {
    CommonLog.lError(`showOverview error: ${e.message}`);
  }
}

function* DebounceViewStep(d, action: tCommonActionType) {
  try {
    switch (action.type) {
      case ORDER.STEP.PREVIOUS:
        yield put({ type: ORDER.STEP.VIEW_PREVIOUS });
        break;
      case ORDER.STEP.NEXT:
        yield put({ type: ORDER.STEP.VIEW_NEXT });
        break;
      default:
        break;
    }
    yield delay(d);
  } catch (e) {
    CommonLog.lError(e);
  }
}

export default function* root(): Saga<void> {
  try {
    yield all([
      takeLatest(ORDER.WORK_ON, workOnOrder),
      takeEvery(ORDER.VIEW, viewOrder),
      takeLeading([ORDER.STEP.PREVIOUS, ORDER.STEP.NEXT], DebounceViewStep, 300)
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* viewOrder({ order }) {
  try {
    yield call(mapping.onOrderView, order);
  } catch (e) {
    CommonLog.lError(e);
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
    CommonLog.Info('Order Finished');
    if (finish) {
      yield call(mapping.onOrderFinish, wOrder);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* doOrder() {
  try {
    while (true) {
      CommonLog.Info('Doing Order...');
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
    CommonLog.lError(e);
  }
}
