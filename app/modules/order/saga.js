// @flow
import {
  call,
  select,
  put,
  all,
  takeEvery,
  takeLeading,
  delay,
  take,
  race
} from 'redux-saga/effects';
import React from 'react';
import type { Saga } from 'redux-saga';
import { orderActions } from './action';
import { workingOrder, orderSteps, doable, viewingOrder } from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { CommonLog } from '../../common/utils';
import type { tStep, tStepArray } from './interface/typeDef';
import type { tCommonActionType } from '../../common/type';
import { orderDetailApi, orderListApi } from '../../api/order';
import { ORDER_STATUS, ORDER } from './constents';
import { bindRushAction } from '../rush/rushHealthz';
import loadingActions from '../loading/action';
import NotifierActions from '../Notifier/action';
import type { updateStateActionType } from './action';
import type { tClsStep } from '../step/Step';
import type { tClsOrder } from './Order';

export default function* root(): Saga<void> {
  try {
    yield all([
      call(bindRushAction.onConnect, orderActions.getList),
      // TODO: should we use takeEvery?
      takeEvery(ORDER.LIST.GET, getOrderList),
      takeEvery(ORDER.DETAIL.GET, getOrderDetail),
      takeEvery(ORDER.WORK_ON, workOnOrder),
      takeEvery(ORDER.VIEW, viewOrder),
      takeEvery(ORDER.NEW, newOrder),
      takeEvery(ORDER.STEP.STATUS, updateStatus),
      takeLeading([ORDER.STEP.PREVIOUS, ORDER.STEP.NEXT], DebounceViewStep, 300)
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* updateStatus(action: updateStateActionType) {
  try {
    const { step } = action;
    if (Reflect.has(step, 'updateStatus')) {
      const st = (step: tClsStep);
      yield call([st, st.updateStatus], action);
    } else {
      throw new Error(
        'updateStatus Error, The Step Without updateStatus Property'
      );
    }
  } catch (e) {
    CommonLog.lError(e, {
      at: 'order.saga.updateStatus'
    });
  }
}

function* newOrder() {
  try {
    yield put(NotifierActions.enqueueSnackbar('Info', '收到新工单'));
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* workOnOrder({ order }) {
  try {
    yield race([
      call(order.run, ORDER_STATUS.WIP),
      take(a => a.type === ORDER.FINISH && a.order === order)
    ]);
  } catch (e) {
    CommonLog.lError(e);
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

function* getOrderDetail({ order }) {
  const rOrder: tClsOrder = (order: tClsOrder);
  try {
    const resp = yield call(orderDetailApi, rOrder.id);
    if (resp.result !== 0) {
      yield put(orderActions.getDetailFail());
    }
  } catch (e) {
    CommonLog.lError(e, {
      at: 'getOrderDetail'
    });
  }
}

function* getOrderList() {
  try {
    yield call(orderListApi);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'getOrderList'
    });
  }
}

function* viewOrder({ order }) {
  // eslint-disable-next-line no-unused-expressions
  (order: tClsOrder);
  try {
    const WIPOrder: tClsOrder = yield select(s => workingOrder(s.order));

    if (WIPOrder === order) {
      // 进行中的工单不显示概览对话框
      return;
    }
    yield put(loadingActions.start());
    yield all([
      call(getOrderDetail, { order }),
      race([take(ORDER.DETAIL.SUCCESS), take(ORDER.DETAIL.FAIL)])
    ]);
    yield put(loadingActions.stop());
    const vOrderSteps: ?tStepArray = yield select(state =>
      orderSteps(viewingOrder(state.order))
    );
    const data =
      (vOrderSteps &&
        vOrderSteps.map((s: tStep, idx) => [
          idx + 1,
          s.name,
          i18n.t(`StepType.${s.type}`),
          s.description
        ])) ||
      [];
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
