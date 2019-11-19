// @flow

import {
  put,
  select,
  takeEvery
} from 'redux-saga/effects';

import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { WORKCENTER_MODE } from './constants';
import switchWorkCenterModeActions from './action';

import { CommonLog } from '../../common/utils';
import { workingOrder } from '../order/selector';
import notifierActions from '../Notifier/action';
import type { tAction } from './interface/typeDef';



function* onSwitchWorkCenterMode(action: tAction): Saga<void> {
  try {
    const WIPworkorder = yield select(s => workingOrder(s.order));
    const { mode } = action;
    if (isNil(WIPworkorder)) {
      // 没有WIP workorder
      yield put(switchWorkCenterModeActions.aSwitchWorkCenterModeValidOK(mode));
    } else {
      yield put(notifierActions.enqueueSnackbar('Error', '当前工位有正在执行的工单,不能切换工单模式'));
    }
  } catch (e) {
    CommonLog.lError(`onSwitchWorkCenterMode Error: ${e.toString()}`);
  }
}


export default function* switchWorkCenterModeRoot(): Saga<void> {
  try {
    yield takeEvery(WORKCENTER_MODE.SWITCH, onSwitchWorkCenterMode);
  } catch (e) {
    CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
  }
}
