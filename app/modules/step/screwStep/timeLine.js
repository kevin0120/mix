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
      resultStatusColor[r.result.toUpperCase()] || resultStatusColor.default,
    footerTitle: r.tool_sn,
    body: `${r.result.toUpperCase()}: wi=${r.wi},mi=${r.mi},ti=${r.ti}`
  }));
