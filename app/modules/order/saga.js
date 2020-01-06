// @flow
import {
  actionChannel,
  all,
  call,
  delay,
  fork,
  join,
  put,
  race,
  select,
  take,
  takeEvery,
  takeLeading
} from 'redux-saga/effects';
import React from 'react';
import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { push } from 'connected-react-router';
import { orderActions } from './action';
import { doable, orderSteps, viewingOrder, workingOrder } from './selector';
import dialogActions from '../dialog/action';
import i18n, { tNS } from '../../i18n';
import Table from '../../components/Table/Table';
import { CommonLog } from '../../common/utils';
import type { tCommonActionType } from '../../common/type';
import { orderDetailByCodeApi, orderListApi, orderReportFinishApi } from '../../api/order';
import { ORDER, ORDER_STATUS } from './constants';
import { bindRushAction } from '../rush/rushHealthz';
import loadingActions from '../loading/action';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { IOrder } from './interface/IOrder';
import notifierActions from '../Notifier/action';
import { bindNewDeviceListener } from '../deviceManager/handlerWSData';
import { sGetWorkCenterMode } from '../workCenterMode/selector';
import type { tWorkCenterMode } from '../workCenterMode/interface/typeDef';
import { stepWorkingNS } from '../../containers/stepWorking/local';
import ClsScanner from '../device/scanner/ClsScanner';
import type { tAnyStatus } from '../step/interface/typeDef';
import type { IWorkable } from '../workable/IWorkable';
import { workModes } from '../workCenterMode/constants';

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
      (input, state) => !workingOrder(state.order),
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
  code,
  config
}: {
  order: IOrder,
  code: string | number
}) {
  try {
    let orderToDo = null;
    if (order) {
      orderToDo = order;
    }
    const workCenterMode: tWorkCenterMode = yield select(s => sGetWorkCenterMode(s));
    if (workCenterMode === workModes.reworkWorkCenterMode) {
      // do nothing when rework
      // yield put(reworkActions.tryRework(orderToDo));
      return;
    }
    const orderState = yield select(s => s.order);
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
    const { list } = orderState;
    const wOrder = list.find(o => o.status === ORDER_STATUS.WIP);
    if (wOrder && wOrder !== order) {
      canWorkOnOrder = false;
    }
    if (canWorkOnOrder) {
      yield put(orderActions.workOn(orderToDo, config));
    } else {
      yield put(notifierActions.enqueueSnackbar('Warn', `无法开始新工单：当前有正在进行的工单`));

    }
  } catch (e) {
    CommonLog.lError(e, { at: 'tryWorkOnOrder' });
  }
}

function* workOnOrder({ order, config }: { order: IOrder }) {
  try {

    yield race([
      call(order.run, ORDER_STATUS.TODO, config),
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

function* getOrderDetail({ code }) {
  try {
    yield put(loadingActions.start());
    const detailResult = yield fork(function* detailResult() {
      yield race([take(ORDER.DETAIL.SUCCESS), take(ORDER.DETAIL.FAIL)]);
    });
    const workCenterCode = yield select(s => s.systemInfo.workcenter);
    yield call(orderDetailByCodeApi, code, workCenterCode);
    // if (resp.result !== 0) {
    //   yield put(orderActions.getDetailFail());
    // }
    yield join(detailResult);
    const orderList = yield select(s => s.order.list);
    const order = orderList.find(o => o.code === code);
    yield put(orderActions.view(order));
    yield put(loadingActions.stop());
  } catch (e) {
    yield put(loadingActions.stop());
    yield put(notifierActions.enqueueSnackbar(
      'Error', `获取工单详情失败（${e.message}）`, {
        at: 'getOrderDetail'
      }
    ));
  }
}

function* getOrderList() {
  try {
    yield call(orderListApi, {
      // TODO: order query args
    });
  } catch (e) {
    yield put(notifierActions.enqueueSnackbar(
      'Error', `获取工单列表失败（${e.message}）`, {
        at: 'getOrderList',
      }
    ));
  }
}

function* tryViewOrder({
  order,
  code
}: {
  order: IOrder,
  code: string | number
}) {
  try {
    yield put(loadingActions.start());
    if (isNil(order) && isNil(code)) {
      yield put(loadingActions.stop());
      return;
    }

    let triggerCode = '';
    if (!isNil(code)) {
      triggerCode = code;
    }
    if (!isNil(order)) {
      triggerCode = order.code;
    }
    // if (!order) {
    //   yield put(loadingActions.stop());
    //   yield put(notifierActions.enqueueSnackbar('Warn', `工单不存在: ${code}`));
    //   return;
    // }
    // yield put(push('/app/working'));
    const vOrder: IOrder = yield select(s => viewingOrder(s.order));
    const WIPOrder: IOrder = yield select(s => workingOrder(s.order));

    if (WIPOrder !== vOrder) {
      yield put(orderActions.clearData(vOrder));
    }
    yield call(getOrderDetail, { code: triggerCode });
  } catch (e) {
    yield put(loadingActions.stop());
    CommonLog.lError(e, { at: 'tryViewOrder', order, code });
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
          tNS(`${s.type}`, stepWorkingNS),
          s.desc
        ])) ||
      [];

    const { workCenterMode } = yield select();
    const isRework = workCenterMode === workModes.reworkWorkCenterMode;
    const oList = yield select(s => s.order.list);
    const wOrder = oList.find(o => o.status === ORDER_STATUS.WIP);
    const showStartButton =
      !isRework &&
      !WIPOrder &&
      (!wOrder || (wOrder === order)) &&
      doable(order);
    yield put(
      dialogActions.dialogShow({
        maxWidth: 'md',
        buttons: [
          {
            label: 'Common.Close',
            color: 'warning'
          },
          showStartButton && {
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
