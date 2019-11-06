// @flow

import { rushSendApi } from './rush';


export function toolEnableApi(toolSN: string, enable: boolean): Promise<any> {
  return rushSendApi('WS_TOOL_ENABLE', {
    tool_sn: toolSN,
    enable
  });
}

export function psetApi(toolSN: string = '', stepID: number, userID: number,
                        pset: string, sequence: number, count: number,
                        total: number, workorderID: string): Promise<any> {
  return rushSendApi('WS_TOOL_PSET', {
    tool_sn: toolSN,
    step_id: stepID,
    user_id: userID,
    total,
    pset,
    sequence,
    count,
    workorder_id: workorderID
  });
}

export function jobApi(toolSN: string = '', stepID: number, userID: number, job: number): Promise<any> {
  return rushSendApi('WS_TOOL_JOB', {
    tool_sn: toolSN,
    step_id: stepID,
    user_id: userID,
    job
  });
}
