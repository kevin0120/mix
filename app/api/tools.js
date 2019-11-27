// @flow

import { rushSendApi } from './rush';


export function toolEnableApi(toolSN: string, ControllerSN: string = '', enable: boolean): Promise<any> {
  return rushSendApi('WS_TOOL_ENABLE', {
    tool_sn: toolSN,
    controller_sn: ControllerSN,
    enable
  });
}

export function psetApi(toolSN: string = '', ControllerSN: string = '', stepID: number, userIDs: Array<number>,
                        pset: string, sequence: number, count: number,
                        total: number, workorderID: number): Promise<any> {
  let psetNum=pset;
  if(typeof pset === 'string'){
    psetNum=parseInt(pset,10);
  }
  return rushSendApi('WS_TOOL_PSET', {
    tool_sn: toolSN,
    controller_sn: ControllerSN,
    step_id: stepID,
    user_id: userIDs[0],
    total,
    pset:psetNum,
    sequence,
    count,
    workorder_id: workorderID
  });
}

export function jobApi(toolSN: string = '', ControllerSN: string = '', stepID: number, userID: number, job: number): Promise<any> {
  return rushSendApi('WS_TOOL_JOB', {
    tool_sn: toolSN,
    controller_sn: ControllerSN,
    step_id: stepID,
    user_id: userID,
    job
  });
}
