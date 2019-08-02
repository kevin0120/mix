// @flow

export const SYSTEM_INIT = 'SYSTEM_INIT';

export const systemInit = (section = 'all') => ({
  type: SYSTEM_INIT,
  section
});
