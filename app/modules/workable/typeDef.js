// @flow

import type { tAnyStatus } from '../step/interface/typeDef';

export type tWorkableData = {
  id?: number,
  code?: string,
  status: tAnyStatus,
  payload: any,
  data: any,
  steps: Array<tWorkableData>
};