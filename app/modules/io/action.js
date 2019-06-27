
export const IO = {
  FUNCTION: 'IO_FUNCTION',
  INIT: 'IO_INIT',
  TEST: 'IO_TEST',
  RESET: 'IO_RESET'
};

export function initIO() {
  return {
    type: IO.INIT
  };
}

export function resetIO(modbusConfig) {
  return {
    type: IO.RESET,
    modbusConfig
  };
}
