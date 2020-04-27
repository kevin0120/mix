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
import i18n from '../../i18n';
import OrderInfoTable from '../../components/OrderInfoTable';
import { CommonLog } from '../../common/utils';
import type { tCommonActionType } from '../../common/type';
import {
  apiOrderStartSimulate,
  getBlockReasonsApi,
  orderDetailByCodeApi,
  orderListApi,
  orderPendingApi,
  orderReportFinishApi,
  orderResumeApi
} from '../../api/order';
import { ORDER, ORDER_STATUS } from './constants';
import { bindRushAction } from '../rush/rushHealthz';
import loadingActions from '../loading/action';
import type { IWorkStep } from '../step/interface/IWorkStep';
import type { IOrder } from './interface/IOrder';
import notifierActions from '../Notifier/action';
import { bindNewDeviceListener } from '../deviceManager/handlerWSData';
import { sGetWorkCenterMode } from '../workCenterMode/selector';
import type { tWorkCenterMode } from '../workCenterMode/interface/typeDef';
import ClsScanner from '../device/scanner/ClsScanner';
import type { tAnyStatus } from '../step/interface/typeDef';
import type { IWorkable } from '../workable/IWorkable';
import { workModes } from '../workCenterMode/constants';
import type { tWorkOnOrderConfig, 工单号 } from './interface/typeDef';
import viewOperationInfo from '../viewOperationInfo';


export default function* root(): Saga<void> {
  try {
    yield takeEvery(ORDER.NEW_SCANNER, onNewScanner);
    yield call(bindNewScanner);
    yield call(getBlockReasons);
    yield all([
      call(bindRushAction.onConnect, orderActions.getList), // 绑定rush连接时需要触发的action
      takeEvery(ORDER.LIST.GET, getOrderList),
      takeEvery(ORDER.DETAIL.GET, getOrderDetail),
      call(watchOrderTrigger),
      takeEvery(ORDER.WORK_ON, workOnOrder),
      takeEvery(a => a.type === ORDER.STEP.STATUS && a.status === ORDER_STATUS.PENDING, handlePending),
      takeEvery(ORDER.VIEW, viewOrder),
      takeEvery(ORDER.TRY_VIEW, tryViewOrder),
      takeEvery(ORDER.REPORT_FINISH, reportFinish),
      takeEvery(ORDER.STEP.STATUS, setStepStatus),
      takeEvery(ORDER.REDO_ORDER, redoOrder),
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
      (input, state) => !workingOrder(state.order) && !state.manual?.working,
      input => orderActions.tryViewCode(input.data)
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'onNewScanner' });
  }
}

function* reportFinish({ order }) {
  try {
    const code = (order: IWorkable)._code;
    const workCenterCode = yield select(s => s.systemInfo.workcenter);
    const dateComplete = new Date();
    const resp = yield call(
      orderReportFinishApi,
      code,
      workCenterCode,
      dateComplete
    );
    if (resp) {
      yield put(notifierActions.enqueueSnackbar('Info', '完工请求完成', {
        resp
      }));
    }
  } catch (e) {
    yield put(notifierActions.enqueueSnackbar('Error', e.message, {
      at: 'reportFinish'
    }));
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
  code: 工单号,
  config: tWorkOnOrderConfig
}) {
  try {
    yield put(loadingActions.start());
    let orderToDo: ?IOrder = null;
    if (order) {
      orderToDo = order;
    }
    const workCenterMode: tWorkCenterMode = yield select(s =>
      sGetWorkCenterMode(s)
    );
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
    if (hasAnotherWorkingOrder(orderState, orderToDo)) {
      throw new Error('无法开始新工单：当前有正在进行的工单');
    }

    yield call(orderStartSimulate, orderToDo.code);
    yield put(loadingActions.stop());
    yield put(orderActions.workOn(orderToDo, config));
  } catch (e) {
    yield put(loadingActions.stop());
    yield put(notifierActions.enqueueSnackbar('Error', e.message));
    CommonLog.lError(e, { at: 'tryWorkOnOrder' });
  }
}

function hasAnotherWorkingOrder(orderState, order: IOrder) {
  if (workingOrder(orderState)) {
    return true;
  }
  const { list } = orderState;
  const wOrder = list.find(o => o.status === ORDER_STATUS.WIP);
  return !!(wOrder && wOrder !== order);
}

function* orderStartSimulate(code: 工单号) {
  try {
    const { users } = yield select();
    const uuids = users.map(u => u.uuid || ''); // todo get users
    const workCenterCode = yield select(s => s.systemInfo.workcenter);
    const { errorMessage } = yield call(
      apiOrderStartSimulate,
      code,
      uuids,
      workCenterCode
    );
    if (errorMessage) {
      throw new Error(errorMessage);
    }
  } catch (e) {
    const { strictOrderSimulate } = yield select(
      s => s.setting.systemSettings
    );
    if (strictOrderSimulate) {
      throw new Error(`产前模拟失败，无法开始作业：${e.message}`);
    }
    yield put(
      notifierActions.enqueueSnackbar('Warn', `产前模拟失败：${e.message}`)
    );
  }
}

function* workOnOrder({
  order,
  config
}: {
  order: IOrder,
  config: tWorkOnOrderConfig
}) {
  try {
    if (order.status === ORDER_STATUS.PENDING) {
      const startTime = new Date();
      const orderCode = order.code;
      const workCenterCode = yield select(s => s.systemInfo.workcenter);
      try {
        yield call(orderResumeApi, startTime, orderCode, workCenterCode);
      } catch (e) {
        const { strictResume } = yield select(s => s.setting.systemSettings);
        if (strictResume) {
          throw e;
        }
      }
    }
    let startStatus = ORDER_STATUS.WIP;
    if (!order.status || order.status === ORDER_STATUS.TODO) {
      startStatus = ORDER_STATUS.TODO;
    }
    yield race([
      call(order.run, startStatus, config),
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
        at: 'getOrderList'
      }
    ));
  }
}

function* tryViewOrder({ order, code }: { order?: IOrder, code?: 工单号 }) {
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
    if (order) {
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

    const { workCenterMode } = yield select();
    const isRework = workCenterMode === workModes.reworkWorkCenterMode;
    const oList = yield select(s => s.order.list);
    const wOrder = oList.find(o => o.status === ORDER_STATUS.WIP);
    const showStartButton =
      !isRework && !WIPOrder && (!wOrder || wOrder === order) && doable(order);
    const showViewTracing = !!([ORDER_STATUS.DONE, ORDER_STATUS.FAIL, ORDER_STATUS.CANCEL].find((s) => s === order?.status));
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
          },
          showViewTracing && {
            label: '生产追溯',
            color: 'info',
            action: viewOperationInfo.action.viewOperationTracing()
          }
        ],
        title: i18n.t('Order.Overview'),
        content: <OrderInfoTable steps={vOrderSteps}/>
      })
    );
  } catch (e) {
    CommonLog.lError(`showOverview error: ${e.message}`);
  }
}

function* getBlockReasons() {
  try {
    const odooUrl = yield select(s => s.setting.page.odooConnection.odooUrl.value);
    const resp = yield call(getBlockReasonsApi, odooUrl);
    if (!resp || !resp.data) {
      // todo : handle no data
      throw new Error('got empty block reason');
    }
    const blockReasons = resp.data.map(r => ({
      name: r.name,
      lossType: r.type
    }));
    yield put(orderActions.setBlockReasonList(blockReasons));
  } catch (e) {
    CommonLog.lError(e, {
      at: 'order getBlockReasons'
    });
  }
}

function* handlePending({ config, step: order }) {
  try {
    if (!order) {
      // todo handle no order
      throw new Error('trying to pending without order');
    }
    let reason = config?.reason;
    if (!reason) {
      reason = {
        lossType: 'availability',
        name: 'Equipment Failure'
      };
    }
    const { lossType: exceptType, name: exceptCode } = reason;

    const PendingTime = new Date();
    const orderCode = order.code;
    const workCenterCode = yield select(s => s.systemInfo.workcenter);
    yield call(orderPendingApi, exceptType, exceptCode, PendingTime, orderCode, workCenterCode);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'order handlePending'
    });
  }
}

function* redoOrder({ order }: { order: IOrder }) {
  try {
    if (!order) {
      throw new Error('没有指定工单');
    }
    // clear order data
    const { steps } = order;
    console.log(steps);
    const clearDataEffects = steps.map(s =>
      call([s, s.reset])
    );
    yield all(clearDataEffects);
    // run order
    yield put(orderActions.tryWorkOn(order));
  } catch (e) {
    yield put(notifierActions.enqueueSnackbar('Error', `工单重新作业失败：${e.message}`, {
      e,
      at: 'redoOrder'
    }));
  }
}
