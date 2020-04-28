// @flow

import { put, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import notifierActions from '../../Notifier/action';
import { CommonLog } from '../../../common/utils';
import type { tDeviceSN } from '../typeDef';
import type { tResult } from '../../step/screwStep/interface/typeDef';
import { getDevice } from '../../deviceManager/devices';

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
  data: Array<{
    tool_sn: tDeviceSN,
    results: Array<tResult>
  }>
};

export function* toolStatusChange(data: tToolStatusData): Saga<void> {
  try {
    if (!data.data) {
      return;
    }
    if (data.type === 'WS_TIGHTENING_TOOL_STATUS') {
      const { tool_sn: toolSN, status, reason } = data.data;
      yield put(
        notifierActions.enqueueSnackbar(
          'Info',
          `拧紧枪状态更新（${toolSN}）： ${String(status)}${
            reason ? `, ${reason}` : ''
          }`
        )
      );
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event tool' });
  }
}

export function* toolNewResults({ data: results }: tToolResultData): Saga<void> {
  try {
    CommonLog.Info('rush result', {
      results: JSON.stringify(results)
    });
    // eslint-disable-next-line no-restricted-syntax
    for (const r of results) {
      const tool = getDevice(r.tool_sn);
      if (tool) {
        yield call(tool.doDispatch, [r]);
      } else {
        CommonLog.lError('invalid tool', {
          sn: r.tool_sn
        });
      }
    }
  } catch (e) {
    CommonLog.lError('onToolResult:', e);
  }
}
