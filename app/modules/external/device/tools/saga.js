// @flow

import { put } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { setNewNotification } from '../../../notification/action';
import ClsScrewTool, { defaultScrewToolDispatcher } from './model';
import type { tRushWebSocketData } from '../../../rush/type';
import { CommonLog } from '../../../../common/utils';

export const staticScrewTool = new ClsScrewTool('G1', '0001');
staticScrewTool.dispatcher = defaultScrewToolDispatcher;


export function* toolStatusChange(data: tRushWebSocketData): Saga<void> {
  try {
    const { tool_sn: toolSN, status, reason } = data.data;
    yield put(
      setNewNotification(
        'Info',
        `拧紧枪状态更新（${toolSN}）： ${status}${reason ? `, ${reason}` : ''}`
      )
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event tool' });
  }
}

export function* toolNewResults(data: tRushWebSocketData): Saga<void> {
  try {
    const { results } = data.data;
    const respAction = staticScrewTool.doDispatch(results);
    if (respAction) {
      yield put(respAction);
    }
  } catch (e) {
    console.error('onToolResult:', e);
  }
}
