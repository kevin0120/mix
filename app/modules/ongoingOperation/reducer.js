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
