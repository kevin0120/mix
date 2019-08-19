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

export function* psetApi(toolSN: string = '', stepID: number, userID: number,
                         pset: number, sequence: number, count: number,
                         total: number): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_TOOL_PSET', {
      tool_sn: toolSN,
      step_id: stepID,
      user_id: userID,
      total,
      pset,
      sequence,
      count
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi', toolSN, stepID, userID, pset, sequence, count
    });
  }
}

export function* jobApi(toolSN: string = '', stepID: number, userID: number, job: number): Saga<void> {
  try {
    return yield call(rushSendApi, 'WS_TOOL_JOB', {
      tool_sn: toolSN,
      step_id: stepID,
      user_id: userID,
      job
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi', toolSN, stepID, userID, job
    });
  }
}
