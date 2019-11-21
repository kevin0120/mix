// red | green | white | yellow
export const ioOutputs = Object.freeze({
  red: 'red',
  green: 'green',
  white: 'white',
  yellow: 'yellow'
});

// resetKey | byPass| modeSelect
export const ioInputs = Object.freeze({
  resetKey: 'resetKey',
  byPass: 'byPass',
  modeSelect: 'modeSelect'
});

// info | warning | error
const { red, green, white, yellow } = ioOutputs;
const singleGroups = {};
Object.values(ioOutputs).forEach(o => {
  singleGroups[o] = [o];
});
export const ioOutputGroups = {
  ...singleGroups,
  info: [white],
  warning: [yellow],
  error: [red],
  ready: [green],
  doing: [white, green]
};

export const IO = {
  ADD_LISTENER: 'IO_ADD_LISTENER',
  SET: 'IO_SET',
  SET_PORT: 'IO_SET_PORT'
};
