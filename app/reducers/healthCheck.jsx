/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import SET_HEALTHZ_CHECK from '../actions/actionTypes';

const lodash = require('lodash');

const defaultHealthChecks = {
  masterpc: {
    displayOrder: 1,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.MasterPC'
  },
  modbus: {
    displayOrder: 21,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.IO'
  },
  rfid: {
    displayOrder: 31,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.RFID'
  },
  odoo: {
    displayOrder: 41,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.Odoo'
  }
};

type actionType = {
  +type: string,
  +section: string,
  +isHealth: boolean
};

export default function healthCheckResults(
  state: object = defaultHealthChecks,
  action: actionType
) {
  switch (action.type) {
    case SET_HEALTHZ_CHECK:
      if (!lodash.has(state, action.section)) return state;
      return {
        ...state,
        [action.section]: { ...[action.section], isHealth: action.isHealth }
      };
    default:
      return state;
  }
}
