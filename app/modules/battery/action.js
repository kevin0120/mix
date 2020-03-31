export const BATTERY = {
  CHECK: 'BATTERY_CHECK',
  CHECK_OK: 'BATTERY_CHECK_OK'
};

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
