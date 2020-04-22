// @flow
import { call, put, race, select, take, takeEvery,fork,cancel } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import type { tAction } from './interface/typeDef';
import { orderActions } from '../order/action';
import dialogActions from '../dialog/action';
import { tNS } from '../../i18n';
import { REWORK_PATTERN, reworkDialogConstants as dia, reworkNS } from './constants';
import { translation as trans } from './local';
import notifierActions from '../Notifier/action';
import actions from './action';
import { workModes } from '../workCenterMode/constants';
import { workingOrder } from '../order/selector';
import { manualResult } from '../manual/saga';
import { getDevice } from '../deviceManager/devices';
import ClsIOModule from '../device/io/ClsIOModule';
import { ioDirection, ioTriggerMode } from '../device/io/constants';
import setting from '../setting/reducer';
import { ioContactApi } from '../../api/io';
import { MATERIAL_STEP } from '../step/materialStep/action';

function* tryRework(action: tAction = {}): Saga<void> {
  try {
    console.log('try rework');
    const { order, step, point } = action;
    let canRework = true;         // check if point can rework
    const { workCenterMode } = yield select();
    const wOrder = yield select(s => workingOrder(s.order));

    if (canRework){
    if (workCenterMode !== workModes.reworkWorkCenterMode) {
      yield put(notifierActions.enqueueSnackbar('Warn', '当前工作模式无法进行返工作业，请先切换至返工模式!'));
      canRework = false;
    } else if (wOrder) {
      yield put(notifierActions.enqueueSnackbar('Error', '当前工位有正在执行的工单,不能切换工单模式'));
      canRework = false;
    } else if (!order) {
      yield put(notifierActions.enqueueSnackbar('Error', '没有指定返工工单!'));
      canRework = false;
    } else if (order && !order.hasUnsuccessWorkStep()) {
      yield put(notifierActions.enqueueSnackbar('Error', '当前工单没有可返工的工步'));
      canRework = false;
    } else if (point && point.noRedo) {
      yield put(notifierActions.enqueueSnackbar('Error', '此拧紧点不具备返工条件!'));
      canRework = false;
    }
    }
    if (canRework) {
      const btnCancel = {
        label: tNS(dia.cancel, reworkNS),
        color: 'warning',
        action: actions.cancelRework()
      };
      const btnConfirm = {
        label: tNS(dia.confirm, reworkNS),
        color: 'info',
        action: actions.doRework(order, step, point)
      };

      const title = point ? tNS(trans.redoSpecScrewPointTitle, reworkNS) :
        tNS(trans.redoSpecScrewPointTitleNoPoint, reworkNS);
      const content = point ? `${tNS(trans.redoSpecScrewPointContent, reworkNS)} ${JSON.stringify(
        {
          Bolt: point?.nut_no,
          Sequence: point?.sequence
        })}` : `${tNS(trans.redoSpecScrewPointContentNoPoint, reworkNS)}`;
      yield put(dialogActions.dialogShow({
        buttons: [btnCancel, btnConfirm],
        title,
        content
      }));
    } else {
      yield put(actions.cancelRework());
    }
  } catch (e) {
    CommonLog.lError(`tryRework Error: ${e.toString()}`);
    yield put(actions.cancelRework());
  }
}

function* doRework(action = {}): Saga<void> {
  try {
    const { order, step, point } = action;
    console.log(order, step, point);
    CommonLog.Debug(`Now Try To Rework Tightening Point: ${point ? JSON.stringify(
      {
        Bolt: point.nut_no,
        Sequence: point.sequence
      }) : 'no point'}`);
    yield put(orderActions.workOn(order, {
      step,
      reworkConfig: {
        point
      }
    })); // 将工单切换到工作状态。
    yield put(actions.cancelRework());
  } catch (e) {
    CommonLog.lError(e);
  }
}

let manual;
let newBool;


export default function* reworkPatternRoot(): Saga<void> {
  try {
    yield takeEvery(REWORK_PATTERN.DO_REWORK, doRework);
  } catch (e) {
    CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
  }
  while (true) {
    try {
      const action = yield take(REWORK_PATTERN.TRY_REWORK);
      const mode = yield select(s => s.workCenterMode);

      // 手动输入拧紧结果的权限控制
      const { adminKey } = yield select(s => s.setting);
      // const confirmIO = getDevice(adminKey.io_sn);
      // if (!(confirmIO instanceof ClsIOModule)) {
      //   throw new Error(`io module (${adminKey.io_sn}) not found`);
      // }
      // const input = confirmIO.getPort(ioDirection.input, adminKey.input);

      try {
      const value=yield call((ioContactApi: Function), adminKey.io_sn);
      newBool = ClsIOModule.bitString2Boolean(value.data?.inputs[adminKey.input])
      } catch (e) {
        CommonLog.lError(`ioContactApi Error: ${e.toString()}`);
      }

      if (mode === workModes.normWorkCenterMode) {
        if (newBool || !adminKey.enable){
          if (manual !==null&& typeof manual !== 'undefined'){
            yield cancel(manual);
          }
          manual =yield fork(manualResult,action);
        } else {
          yield put(notifierActions.enqueueSnackbar('Warn', '当前工作状态无法进行手动输入结果!请检查钥匙开关是否打开'));
        }

      } else {
        yield race([
          call(tryRework, action),
          take(REWORK_PATTERN.CANCEL_REWORK)
        ]);
      }
    } catch (e) {
      CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
    }
  }
}
