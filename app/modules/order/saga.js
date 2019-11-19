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
  join,
  actionChannel
} from 'redux-saga/effects';
import React from 'react';
import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { push } from 'connected-react-router';
import { orderActions } from './action';
import { workingOrder, orderSteps, doable, viewingOrder } from './selector';
import dialogActions from '../dialog/action';
import i18n from '../../i18n';
import Table from '../../components/Table/Table';
import { CommonLog } from '../../common/utils';
import type { tCommonActionType } from '../../common/type';
import {
  orderDetailApi,
  orderDetailByCodeApi,
  orderListApi,
  orderReportFinishApi
} from '../../api/order';
import { ORDER, ORDER_STATUS } from './constants';
import { bindRushAction } from '../rush/rushHealthz';
import loadingActions from '../loading/action';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { IOrder } from './interface/IOrder';
import notifierActions from '../Notifier/action';
import { bindNewDeviceListener } from '../deviceManager/handlerWSData';
import {sGetWorkCenterMode} from '../workCenterMode/selector';
import { translation as trans } from '../../components/NavBar/local';
import type {tWorkCenterMode} from  '../workCenterMode/interface/typeDef'

import ClsScanner from '../device/scanner/ClsScanner';
import type { tAnyStatus } from '../step/interface/typeDef';
import type { IWorkable } from '../workable/IWorkable';

export default function* root(): Saga<void> {
  try {
    yield takeEvery(ORDER.NEW_SCANNER, onNewScanner);
    yield call(bindNewScanner);
    yield all([
      call(bindRushAction.onConnect, orderActions.getList), // 绑定rush连接时需要触发的action
      takeEvery(ORDER.LIST.GET, getOrderList),
      takeEvery(ORDER.DETAIL.GET, getOrderDetail),
      call(watchOrderTrigger),
      takeEvery(ORDER.WORK_ON, workOnOrder),
      takeEvery(ORDER.VIEW, viewOrder),
      takeEvery(ORDER.TRY_VIEW, tryViewOrder),
      takeEvery(ORDER.REPORT_FINISH, reportFinish),
      takeEvery(ORDER.STEP.STATUS, setStepStatus),
      takeLeading([ORDER.STEP.PREVIOUS, ORDER.STEP.NEXT], DebounceViewStep, 300)
    ]);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* setStepStatus({
  step,
  status
}: {
  step: IWorkable,
  status: tAnyStatus
}) {
  if (!step || !status) {
    return;
  }
  try {
    yield call([step, step.updateStatus], { status });
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindNewScanner() {
  try {
    yield call(
      bindNewDeviceListener,
      d => d instanceof ClsScanner,
      s => orderActions.newScanner(s)
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'bindNewScanner' });
  }
}

// 扫码触发工单
function onNewScanner({ scanner }) {
  try {
    // TODO: filter scanner input
    scanner.Enable();
    scanner.addListener(
      () => true,
      input => orderActions.tryViewCode(input.data)
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'onNewScanner' });
  }
}

function* reportFinish({
  code,
  trackCode,
  productCode,
  workCenterCode,
  dateComplete,
  operation
}) {
  try {
    yield call(
      orderReportFinishApi,
      code,
      trackCode,
      productCode,
      workCenterCode,
      dateComplete,
      operation
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'reportFinish' });
  }
}

function* watchOrderTrigger() {
  try {
    const triggerChannel = yield actionChannel(ORDER.TRY_WORK_ON);
    while (true) {
      const action = yield take(triggerChannel);
      yield call(tryWorkOnOrder, action);
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'orderTrigger' });
  }
}

function* tryWorkOnOrder({
  order,
  code
}: {
  order: IOrder,
  code: string | number
}) {
  try {
    const orderState = yield select(s => s.order);
    const workCenterMode: tWorkCenterMode = yield select(s => sGetWorkCenterMode(s));
    if (workCenterMode === trans.reworkWorkCenterMode && !order.hasFailWorkStep()) {
      // 在返工模式下，但此工单并没有失败的工步
      yield put(notifierActions.enqueueSnackbar('Error', '当前工单没有可返工的工步'));
      return;
    }
    let orderToDo = null;
    if (order) {
      orderToDo = order;
    }
    if (code) {
      const { list: orderList } = orderState;
      orderToDo = orderList.find(o => o.code === code);
    }
    if (!orderToDo) {
      return;
    }

    let canWorkOnOrder = true;
    if (workingOrder(orderState)) {
      canWorkOnOrder = false;
    }
    if (canWorkOnOrder) {
      yield put(orderActions.workOn((orderToDo: IOrder)));
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'tryWorkOnOrder' });
  }
}

function* workOnOrder({ order }: { order: IOrder }) {
  try {
    yield race([
      call(order.run, ORDER_STATUS.TODO),
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
// TODO: auto set wip order working
function* getOrderDetail({ order }) {
  try {
    yield put(loadingActions.start());
    const detailResult = yield fork(function* detailResult() {
      yield race([take(ORDER.DETAIL.SUCCESS), take(ORDER.DETAIL.FAIL)]);
    });
    yield call(orderDetailApi, order.id);
    // if (resp.result !== 0) {
    //   yield put(orderActions.getDetailFail());
    // }
    yield join(detailResult);
    yield put(loadingActions.stop());
  } catch (e) {
    yield put(loadingActions.stop());
    CommonLog.lError(e, {
      at: 'getOrderDetail',
      code: order.code
    });
  }
}

function* getOrderList() {
  try {
    yield call(orderListApi, {
      // TODO: order query args
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'getOrderList'
    });
  }
}

function* tryViewOrder({
  order: orderRec,
  code
}: {
  order: IOrder,
  code: string | number
}) {
  try {
    yield put(loadingActions.start());

    const orderList = yield select(s => s.order.list);
    let order = orderRec;

    if (isNil(order) && !isNil(code)) {
      order = orderList.find(o => o.code === code);
    }

    if (!order) {
      const workcenterCode = yield select(s => s.systemInfo.workcenter);
      const { data } = yield call(orderDetailByCodeApi, code, workcenterCode);
      order = data;
    }

    // TODO: check no-trigger conditions
    if (!order) {
      yield put(loadingActions.stop());
      yield put(notifierActions.enqueueSnackbar('Warn', `工单不存在: ${code}`));
      return;
    }
    yield put(push('/app/working'));

    yield call(getOrderDetail, { order });

    yield put(orderActions.view(order));
  } catch (e) {
    yield put(loadingActions.stop());
    CommonLog.lError(e, { at: 'tryViewOrder', orderRec, code });
  }
}

function* viewOrder({ order }: { order: IOrder }) {
  try {
    yield put(push('/app/working'));

    const WIPOrder: IOrder = yield select(s => workingOrder(s.order));

    if (WIPOrder === order) {
      // 进行中的工单不显示概览对话框
      return;
    }

    const vOrderSteps: ?Array<IWorkStep> = yield select(state =>
      orderSteps(viewingOrder(state.order))
    );

    const data =
      (vOrderSteps &&
        vOrderSteps.map((s: IWorkStep, idx) => [
          idx + 1,
          s.code,
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
