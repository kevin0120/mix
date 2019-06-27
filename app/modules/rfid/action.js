
export const RFID = {
  INIT: 'RFID_INIT',
  TOGGLE_STATUS: 'RFID_TOGGLE_STATUS'
};

export function toggleRFID() {
  return {
    type: RFID.TOGGLE_STATUS
  };
}
