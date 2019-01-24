import { BATTERY } from './actionTypes';


export function batteryCheck() {
  return {
    type: BATTERY.CHECK
  };
}

export function batteryCheckOK(percentage) {
  return {
    type: BATTERY.CHECK_OK,
    percentage
  };
}
