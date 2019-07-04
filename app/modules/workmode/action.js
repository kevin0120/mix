// @flow

export const WORK_MODE = {
  SWITCH_WM: 'SWITCH_WORK_MODE',
  SWITCH_CM: 'SWITCH_CONTROL_MODE'
};

export function switchWorkMode(mode: string) {
  return {
    type: WORK_MODE.SWITCH_WM,
    mode
  };
}
