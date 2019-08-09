// @flow
import { call, delay, race } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../common/utils';
import rushSendMessage from '../modules/rush/sendMessage';

const timeout = 10000;

export function* psetApi(toolSN: string = '', stepID: number, userID: number,
                         pset: number, sequence: number, count: number): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_PSET',
        data: {
          tool_sn: toolSN,
          step_id: stepID,
          user_id: userID,
          pset,
          sequence,
          count
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'psetApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi', toolSN, stepID, userID, pset, sequence, count
    });
  }
}

export function* jobApi(toolSN: string = '', stepID: number, userID: number, job: number): Saga<void> {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_JOB',
        data: {
          tool_sn: toolSN,
          step_id: stepID,
          user_id: userID,
          job
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return {
        result: -1,
        msg: 'jobApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi', toolSN, stepID, userID, job
    });
  }
}



