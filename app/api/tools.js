// @flow

import { rushSendApi } from './rush';


export function toolEnableApi(toolSN: string, ControllerSN: string = '', enable: boolean): Promise<any> {
  return rushSendApi('WS_TOOL_ENABLE', {
    tool_sn: toolSN,
    controller_sn: ControllerSN,
    enable
  });
}

export function psetApi(toolSN: string = '', ControllerSN: string = '', stepCode: number, userIDs: Array<number>,
  pset: string, sequence: number,
  total: number, workorderCode: number, scannerCode: string = '', batch = 0, count = 1): Promise<any> {
  let psetNum = pset;
  if (typeof pset === 'string') {
    psetNum = parseInt(pset, 10);
  }
  return rushSendApi('WS_TOOL_PSET', {
    tool_sn: toolSN,
    controller_sn: ControllerSN,
    workstep_code: stepCode,
    user_id: userIDs[0],
    total,
    pset: psetNum,
    sequence,
    workorder_code: workorderCode,
    scanner_code: scannerCode,
    batch,
    count
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

export function getPestListApi(toolSN, ControllerSN): Promise<any> {
  return rushSendApi('WS_TOOL_PSET_LIST', {
    controller_sn: ControllerSN,
    tool_sn: toolSN
  });
}

export function apiManualSetResult({
  tool_sn,
  controller_sn,
  measure_result,
  measure_torque,
  measure_angle,
  count
}) {
  return rushSendApi('WS_TOOL_RESULT_SET', {
    tool_sn,
    controller_sn,
    measure_result,
    measure_torque,
    measure_angle,
    count
  });
}
