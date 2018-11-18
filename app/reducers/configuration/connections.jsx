/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// flow

import { CONNECTION } from '../../actions/actionTypes';

const defaultConnInfo = {
  masterpc: '',
  rfid: '',
  controllers: [],
  io: '',
  workcenterCode: ''
};

type actionType = {
  +type: string,
  +data: object
};

export default function connections(
  state: object = defaultConnInfo,
  action: actionType
) {
  switch (action.type) {
    case CONNECTION.FETCH_OK:{
      const { masterpc, rfid, io, controllers, info } = action.data;
      return {
        ...state,
        masterpc: masterpc.connection,
        rfid: rfid.connection,
        io: io.connection,
        workcenterCode: info.workcenter_code,
        controllers
      };
    }
    default:
      return state;
  }
}

