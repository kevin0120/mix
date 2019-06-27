/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { ONGOING_OPERATION } from './action';

const defaultOngoingOperation = {
  vin: '',
  model: '',
  lnr: '',
  knr: '',
  long_pin: ''
};

type actionType = {
  +type: string,
  +data: object
};

export default function ongoingOperation(
  state: object = defaultOngoingOperation,
  action: actionType
) {
  switch (action.type) {
    case ONGOING_OPERATION.FETCH_OK: {
      const { vin, model, lnr, knr, long_pin } = action.data;
      return { vin, model, lnr, knr, long_pin };
    }
    case ONGOING_OPERATION.CLEAN: {
      return defaultOngoingOperation;
    }
    default:
      return state;
  }
}
