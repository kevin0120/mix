// @flow
import type {
  tResult,
  tScrewStepData
} from './interface/typeDef';

const resultStatusColor = {
  NOK: 'danger',
  OK: 'success',
  LSN: 'info',
  AK2: 'info',
  default: 'info'
};

export const reduceResult2TimeLine = (results: Array<tResult>) =>
  (d: tScrewStepData): tScrewStepData => ({
    ...d,
    timeLine: [
      ...results.map(r => ({
        title: r.batch,
        color: resultStatusColor[r.result] || resultStatusColor.default,
        footerTitle: r.toolSN,
        body: `${r.result}: wi=${r.wi},mi=${r.mi},ti=${r.ti}`
      })),
      ...(d.timeLine || [])
    ]
  });
