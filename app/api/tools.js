// @flow

import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';


export function* toolEnableApi(toolSN: string, enable: boolean): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_TOOL_ENABLE', {
      tool_sn: toolSN,
      enable
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'toolEnableApi', toolSN, enable
    });
  }
}