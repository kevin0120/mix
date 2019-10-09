// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import ClsController from './model';
import { CommonLog } from '../../../../common/utils';
import type { tRushWebSocketData } from '../../../rush/type';

// eslint-disable-next-line prefer-const
const controller = new ClsController('controller', 'Dummy Serial Number');

export default function* controllerNewData(
  data: tRushWebSocketData
): Saga<void> {
  try {
    yield call(controller.doDispatch, data);
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event controller' });
  }
}
