// @flow

import { SYSTEM_INIT } from './actionTypes';

export const systemInit = (
  baseUrl: string,
  hmiSN: string,
  dispatch,
  getState
) => ({
  type: SYSTEM_INIT,
  baseUrl,
  hmiSN,
  dispatch,
  getState
});
