// @flow

import {
  put,
  take,
  call,
  select,
  takeEvery
} from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { doPoint } from '../step/screwStep/ScrewStep';
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

function* tryRework(action: tAction = {}): Saga<void> {
  try {
    console.log('try rework');
    const { order, step, point } = action;
    let canRework = true;         // check if point can rework
    const { workCenterMode } = yield select();
    if (workCenterMode !== workModes.reworkWorkCenterMode) {
      yield put(notifierActions.enqueueSnackbar('Warn', '当前工作模式无法进行返工作业，请先切换至返工模式!'));
      canRework = false;
    }
    if (!order) {
      yield put(notifierActions.enqueueSnackbar('Error', '没有指定返工工单!'));
      canRework = false;
    }
    if (order && !order.hasFailWorkStep()) {
      yield put(notifierActions.enqueueSnackbar('Error', '当前工单没有可返工的工步'));
      canRework = false;
    }
    if (point && !point.canRedo) {
      yield put(notifierActions.enqueueSnackbar('Error', '此拧紧点不具备返工条件!'));
      canRework = false;
    }

    if (canRework) {
      const btnConfirm = {
        label: tNS(dia.confirm, reworkNS),
        color: 'info',
        action: actions.doRework(order, step, point)
      };
      const btnCancel = {
        label: tNS(dia.cancel, reworkNS),
        color: 'warning',
        action: actions.cancelRework()
      };
      yield put(dialogActions.dialogShow({
        buttons: [btnCancel, btnConfirm],
        title: tNS(trans.redoSpecScrewPointTitle, reworkNS),
        content: `${tNS(trans.redoSpecScrewPointContent, reworkNS)} ${JSON.stringify(
          {
            Bolt: point.nut_no,
            Sequence: point.sequence
          })}`
      }));
    } else {
      yield put(actions.cancelRework());
    }
    // if (!isNil(extra)) {
    //   const { point, step } = extra;
    //
    // } else {
    //   yield put(notifierActions.enqueueSnackbar('Error', '当前工位有正在执行的工单,不能切换工单模式'));
    // }
  } catch (e) {
    CommonLog.lError(`onSwitchWorkCenterMode Error: ${e.toString()}`);
    yield put(actions.cancelRework());
  }
}

function* doRework(action = {}): Saga<void> {
  try {
    const { order, step, point } = action;
    console.log(order, step, point);
    CommonLog.Debug(`Now Try To Rework Tightening Point: ${JSON.stringify(
      {
        Bolt: point.nut_no,
        Sequence: point.sequence
      })}`);
    yield put(orderActions.workOn(order, {
      step,
      reworkConfig: {
        point
      }
    })); // 将工单切换到工作状态。
  } catch (e) {
    CommonLog.lError(e);
  }
}


export default function* reworkPatternRoot(): Saga<void> {
  try {
    yield takeEvery(REWORK_PATTERN.DO_REWORK, doRework);
  } catch (e) {
    CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
  }
  while (true) {
    try {
      const action = yield take(REWORK_PATTERN.TRY_REWORK);
      yield call(tryRework, action);
    } catch (e) {
      CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
    }
  }
}

