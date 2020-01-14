// @flow
import type { tResult } from './interface/typeDef';
import type { tTimeLineObj } from '../../../components/WorkPageTimeline/typeDef';

const resultStatusColor = {
  NOK: 'danger',
  OK: 'success',
  LSN: 'info',
  AK2: 'info',
  default: 'info'
};

// eslint-disable-next-line import/prefer-default-export
export const result2TimeLine = (results: Array<tResult>): Array<tTimeLineObj> =>
  results.map(r => ({
    title: r.batch,
    color:
      resultStatusColor[r.measure_result.toUpperCase()] || resultStatusColor.default,
    footerTitle: r.tool_sn,
    body: `${r.measure_result.toUpperCase()}: torque=${r.measure_torque},angle=${r.measure_angle},time=${r.measure_time}`
  }));
