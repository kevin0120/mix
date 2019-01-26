import { RFID } from './actionTypes';

export function toggleRFID() {
  return {
    type: RFID.TOGGLE_STATUS
  };
}
