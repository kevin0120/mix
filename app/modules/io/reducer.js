import { IO, ioOutputs, ioInputs, ioOutputGroups } from './constants';
import { genReducers } from '../util';
import { ioDirection } from '../device/io/constants';

const reducers = {
  [IO.SET_PORT]: (state, { output, portIdx }) => ({
    ...state,
    [output]: portIdx
  }),
  [IO.SET_MODULE]: (state, { io }) => ({
    ...state,
    ioModule: io
  })
};


const { red, green, white, yellow } = ioOutputs;
const { resetKey, byPass, modeSelect } = ioInputs;
const initState = {
  ioPorts: {
    [ioDirection.output]: {
      [red]: 0,
      [green]: 1,
      [white]: 2,
      [yellow]: 3
    },
    [ioDirection.input]: {
      [resetKey]: 0,
      [byPass]: 1,
      [modeSelect]: 2
    }
  },
  ioOutputGroups,
  ioModule: null
};

export default genReducers(reducers, initState);
