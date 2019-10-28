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
  race,
  fork,
  join
} from 'redux-saga/effects';
import React from 'react';
import type { Saga } from 'redux-saga';
import { orderActions } from './action';
import { workingOrder, orderSteps, doable, viewingOrder } from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { CommonLog } from '../../common/utils';
import type { tCommonActionType } from '../../common/type';
import { orderDetailApi, orderListApi } from '../../api/order';
import { ORDER } from './constants';
import { bindRushAction } from '../rush/rushHealthz';
import loadingActions from '../loading/action';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { IOrder } from './interface/IOrder';

export default function* root(): Saga<void> {
  try {
    yield all([
      call(bindRushAction.onConnect, orderActions.getList), // 绑定rush连接时需要触发的action
      takeEvery(ORDER.LIST.GET, getOrderList),
      takeEvery(ORDER.DETAIL.GET, getOrderDetail),
      takeEvery(ORDER.TRY_WORK_ON, tryWorkOnOrder),
      takeEvery(ORDER.WORK_ON, workOnOrder),
      takeEvery(ORDER.VIEW, viewOrder),
      takeLeading([ORDER.STEP.PREVIOUS, ORDER.STEP.NEXT], DebounceViewStep, 300)
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}

// TODO: 扫码触发工单

// TODO: 工单持久化

// TODO: 开工、报工接口

function* tryWorkOnOrder({ order }: { order: IOrder }) {
  try {
    // TODO: judge can work on
    let canWorkOnOrder = true;
    const orderState = yield select(s => s.order);
    if (workingOrder(orderState)) {
      canWorkOnOrder = false;
    }
    if (!order) {
      canWorkOnOrder = false;
    }
    if (canWorkOnOrder) {
      yield put(orderActions.workOn(order));
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'tryWorkOnOrder' });
  }
}

function* workOnOrder({ order }: { order: IOrder }) {
  try {
    yield race([
      call(order.run),
      take(a => a.type === ORDER.FINISH && a.order === order)
    ]);
    yield put(orderActions.orderDidFinish());
  } catch (e) {
    CommonLog.lError(e, { at: 'workOnOrder' });
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
  const rOrder: IOrder = (order: IOrder);
  try {
    yield put(loadingActions.start());
    const detailResult = yield fork(function* detailResult() {
      yield race([take(ORDER.DETAIL.SUCCESS), take(ORDER.DETAIL.FAIL)]);
    });
    const resp = yield call(orderDetailApi, rOrder.id);
    if (resp.result !== 0) {
      yield put(orderActions.getDetailFail());
    }
    yield join(detailResult);
    yield put(loadingActions.stop());
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

function* viewOrder({ order }: { order: IOrder }) {
  try {
    const WIPOrder: IOrder = yield select(s => workingOrder(s.order));

    if (WIPOrder === order) {
      // 进行中的工单不显示概览对话框
      return;
    }
    yield call(getOrderDetail, { order });

    const vOrderSteps: ?Array<IWorkStep> = yield select(state =>
      orderSteps(viewingOrder(state.order))
    );

    const data =
      (vOrderSteps &&
        vOrderSteps.map((s: IWorkStep, idx) => [
          idx + 1,
          s.name,
          i18n.t(`StepType.${s.type}`),
          s.desc
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
              action: orderActions.tryWorkOn(order)
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
