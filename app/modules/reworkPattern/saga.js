// @flow

import {
  put,
  take,
  call,
  select
} from 'redux-saga/effects';

import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { REWORK_PATTERN } from './constants';
import {doPoint} from "../step/screwStep/ScrewStep";
import { CommonLog } from '../../common/utils';
import notifierActions from '../Notifier/action';
import type { tAction } from './interface/typeDef';
import { orderActions } from "../order/action";
import { viewingOrder } from '../order/selector';


function* onRework(action: tAction): Saga<void> {
  try {
    const { extra } = action;
    
    if (!isNil(extra)) {
      const {point, step} = extra;
      CommonLog.Debug(`Now Try To Rework Tightening Point: ${point.toString()}`);
      const viewOrder = yield select(state => viewingOrder(state.order));
      yield put(orderActions.workOn(viewOrder)); // 将工单切换到工作状态。
      yield call([step,doPoint],[point], false, orderActions);
    } else {
      yield put(notifierActions.enqueueSnackbar('Error', '当前工位有正在执行的工单,不能切换工单模式'));
    }
  } catch (e) {
    CommonLog.lError(`onSwitchWorkCenterMode Error: ${e.toString()}`);
  }
}


export default function* reworkPatternRoot(): Saga<void> {
  while (true) {
    try {
      const action = yield take([REWORK_PATTERN.SPEC_SCREW, REWORK_PATTERN.SCREW_STEP]);
      yield call(onRework, action);
    } catch (e) {
      CommonLog.lError(`switchWorkCenterModeRoot Error: ${e.toString()}`);
    }
  }
}
  
