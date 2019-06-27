// @flow

import { HEALTHZ_CHECK } from './action';

import configs from '../../shared/config';

const lodash = require('lodash');

const defaultHealthChecks = {
  masterpc: {
    displayOrder: 1,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.MasterPC',
    enable: true
  },
  controller: {
    displayOrder: 11,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.Controller',
    enable: true
  },
  modbus: {
    displayOrder: 21,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.IO',
    enable: configs.systemSettings.modbusEnable
  },
  rfid: {
    displayOrder: 31,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.RFID',
    enable: configs.systemSettings.rfidEnabled
  },
  andon: {
    displayOrder: 41,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.Andon',
    enable: configs.systemSettings.andonEnable
  },
  odoo: {
    displayOrder: 51,
    optionInfo: '',
    isHealth: false,
    displayTitle: 'HealthCheck.Odoo',
    enable: true
  }
};

type actionType = {
  +type: string,
  +section: string,
  +isHealth: boolean
};

export default function healthCheckResults(
  state: any = defaultHealthChecks,
  action: actionType
) {
  switch (action.type) {
    case HEALTHZ_CHECK.SET:
      if (!lodash.has(state, action.section)) return state;
      return {
        ...state,
        [action.section]: {
          ...state[action.section],
          isHealth: action.isHealth
        }
      };
    default:
      return state;
  }
}
