// flow

import { CONNECTION } from './action';

import configs from '../../shared/config';

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
      const {
        masterpc,
        rfid,
        io,
        controllers,
        aiis,
        workcenterCode,
        rework_workcenter
      } = action.data;
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
