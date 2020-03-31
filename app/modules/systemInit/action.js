// @flow

export const SYSTEM_INIT = 'SYSTEM_INIT';

export const systemInit = (section: string = 'all') => ({
  type: SYSTEM_INIT,
  section
});
