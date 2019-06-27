/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

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
