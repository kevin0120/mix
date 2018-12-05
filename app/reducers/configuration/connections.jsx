/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// flow

import { CONNECTION } from '../../actions/actionTypes';

import configs from '../../shared/config/index';

const lodash = require('lodash');

const defaultConnInfo = configs.system.connections;

type actionType = {
  +type: string,
  +data: object
};

export default function connections(
  state: object = defaultConnInfo,
  action: actionType
) {
  switch (action.type) {
    case CONNECTION.FETCH_OK: {
      const { masterpc, rfid, io, controllers, info } = action.data;
      return {
        ...state,
        masterpc: masterpc.connection ? masterpc.connection : '',
        rfid: rfid.connection ? rfid.connection : '',
        io: io.connection ? io.connection : '',
        workcenterCode: info.workcenter_code ? info.workcenter_code : '',
        rework_workcenter: info.qc_workcenter ? info.qc_workcenter : '',
        controllers: lodash.isArray(controllers) ? controllers : []
      };
    }
    case CONNECTION.MANUAL_MODIFICATION: {
      const { masterpc, rfid, io, controllers, aiis,workcenterCode, rework_workcenter} = action.data;
      return {
        ...state,
        masterpc,
        aiis,
        rfid,
        io,
        controllers,
        workcenterCode,
        rework_workcenter
      };
    }
    default:
      return state;
  }
}
