// @flow

import {
  put,
  take,
  call
} from 'redux-saga/effects';

import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { REWORK_PATTERN } from './constants';
import switchWorkCenterModeActions from './action';

import { CommonLog } from '../../common/utils';
import notifierActions from '../Notifier/action';
import type { tAction } from './interface/typeDef';


function* onRework(action: tAction): Saga<void> {
  try {
    const { extra } = action;
    if (!isNil(extra)) {
      yield put(switchWorkCenterModeActions.aSwitchWorkCenterModeValidOK(mode));
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
  
