// @flow

export const HEALTHZ_CHECK = {
  SET: 'HEALTHZ_CHECK.SET',
  START: 'HEALTHZ_CHECK.START',
  RESTART: 'HEALTHZ_CHECK.RESTART'
};

export function startHealthzCheck(url = null, controllers = null) {
  return {
    type: HEALTHZ_CHECK.START,
    url,
    controllers
  };
}

export function restartHealthzCheck() {
  return {
    type: HEALTHZ_CHECK.RESTART
  };
}

export function setHealthzCheck(section, isHealth) {
  return {
    type: HEALTHZ_CHECK.SET,
    section,
    isHealth
  };
}
