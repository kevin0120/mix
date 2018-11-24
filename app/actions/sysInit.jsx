// @flow

import { SYSTEM_INIT } from './actionTypes';

export const systemInit = (baseUrl: string, aiisUrl: string, hmiSN: string, dispatch) => ({
  type: SYSTEM_INIT,
  baseUrl,
  aiisUrl,
  hmiSN,
  dispatch
});
