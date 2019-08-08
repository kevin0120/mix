import { call, delay, race } from 'redux-saga/effects';
import { CommonLog } from '../common/utils';
import rushSendMessage from '../modules/rush/sendMessage';

const timeout = 10000;

export function* psetApi(tool_sn = '', step_id, user_id, pset, sequence, count) {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_PSET',
        data: {
          tool_sn,
          step_id,
          user_id,
          pset,
          sequence,
          count
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return{
        result:-1,
        msg:'psetApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'psetApi', tool_sn, step_id, user_id, pset, sequence, count
    });
  }
}

export function* jobApi(tool_sn, step_id, user_id, job) {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_JOB',
        data: {
          tool_sn,
          step_id,
          user_id,
          job
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return{
        result:-1,
        msg:'jobApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi', tool_sn, step_id, user_id, job
    });
  }
}

export function* toolEnableApi(tool_sn, enable) {
  try {
    const { resp, timeout: tOut } = yield race({
      resp: call(rushSendMessage, {
        type: 'WS_TOOL_ENABLE',
        data: {
          tool_sn, enable
        }
      }),
      timeout: delay(timeout)
    });
    if (tOut) {
      return{
        result:-1,
        msg:'toolEnableApi timeout'
      };
    }
    const { data } = resp;
    return data;
  } catch (e) {
    CommonLog.lError(e, {
      at: 'toolEnableApi', tool_sn, enable
    });
  }
}

