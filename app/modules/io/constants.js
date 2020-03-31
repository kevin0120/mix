// red | green | white | yellow
export const ioOutputs = Object.freeze({
  red: 'red',
  green: 'green',
  white: 'white',
  yellow: 'yellow',
  beep: 'beep',
  unlock: 'unlock'
});

// resetKey | byPass| modeSelect | unlock
export const ioInputs = Object.freeze({
  resetKey: 'resetKey',
  byPass: 'byPass',
  modeSelect: 'modeSelect'
});

// info | warning | error
const { red, green, white, yellow } = ioOutputs;
const singleGroups = {};
Object.values(ioOutputs).forEach(o => {
  singleGroups[o] = { on: [o] };
});
export const ioOutputGroups = {
  ...singleGroups,
  info: {
    on: [white],
    off: [red, green, yellow]
  },
  warning: {
    on: [yellow],
    off: [red, green, white]
  },
  error: {
    on: [red],
    off: [yellow, green, white]
  },
  ready: {
    on: [green],
    off: [yellow, red, white]
  },
  doing: {
    on: [white, green],
    off: [yellow, red]
  }
};
