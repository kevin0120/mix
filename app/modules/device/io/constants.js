// @flow

export const ioDirection = Object.freeze({
  input: 'input',
  output: 'output'
});

const IO_FUNCTION = {
  IN: {
    RESET: 'RESET',
    BYPASS: 'BYPASS',
    MODE_SELECT: 'MODE_SELECT'
  },
  OUT: {
    LED_WHITE: 'LED_WHITE',
    LED_YELLOW: 'LED_YELLOW',
    LED_GREEN: 'LED_GREEN',
    LED_RED: 'LED_RED',
    BEEP: 'BEEP'
  }
};

export const IO_WS_TYPES = Object.freeze({
  STATUS: 'WS_IO_STATUS',
  CONTACT: 'WS_IO_CONTACT',
  SET: 'WS_IO_SET'
});

export const ioTriggerMode = {
  rising: 'rising',
  falling: 'falling',
  // high:'high',
  // low:'low',
  change: 'change'
};

export default IO_FUNCTION;
