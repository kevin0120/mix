// @flow

import { put } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { setNewNotification } from '../../../notification/action';
import ClsScrewTool, { defaultScrewToolDispatcher } from './model';
import { CommonLog } from '../../../../common/utils';
import type { tDeviceSN } from '../Device';
import type { tResult } from '../../../step/screwStep/model';

// export const staticScrewTool = new ClsScrewTool('G1', '0001');
// staticScrewTool.dispatcher = defaultScrewToolDispatcher;

type tToolStatusData = {
  type: string,
  data: {
    tool_sn: tDeviceSN,
    status: boolean,
    reason: string
  }
};

type tToolResultData = {
  type: string,
  data: {
    tool_sn: tDeviceSN,
    results: Array<tResult>
  }
};

export function* toolStatusChange(data: tToolStatusData): Saga<void> {
  try {
    const { tool_sn: toolSN, status, reason } = data.data;
    yield put(
      setNewNotification(
        'Info',
        `拧紧枪状态更新（${toolSN}）： ${status.toString()}${reason ? `, ${reason}` : ''}`
      )
    );
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event tool' });
  }
}

export function* toolNewResults(data: tToolResultData): Saga<void> {
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
