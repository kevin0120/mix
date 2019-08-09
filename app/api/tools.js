// @flow

import { call, delay, race } from '@redux-saga/core/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../common/utils';
import rushSendMessage from '../modules/rush/sendMessage';


const timeout = 3000;

export function* toolEnableApi(toolSN: string, enable: boolean): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_ENABLE',
        data: {
          tool_sn: toolSN,
          enable
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'toolEnableApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'toolEnableApi', toolSN, enable
    });
  }
}