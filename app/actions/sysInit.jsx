// @flow

import { SYSTEM_INIT } from './actionTypes';

export const systemInit = (section='all') => ({
  type: SYSTEM_INIT,
  section
});
