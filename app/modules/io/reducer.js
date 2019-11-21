import { IO, ioOutputs, ioInputs, ioOutputGroups } from './constants';
import { genReducers } from '../util';

const reducers = {
  [IO.SET_PORT]: (state, { output, portIdx }) => ({
    ...state,
    [output]: portIdx
  })
};


const { red, green, white, yellow } = ioOutputs;
const { resetKey, byPass, modeSelect } = ioInputs;
const initState = {
  ioPorts: {
    out: {
      [red]: 0,
      [green]: 1,
      [white]: 2,
      [yellow]: 3
    },
    in: {
      [resetKey]: 0,
      [byPass]: 1,
      [modeSelect]: 2
    }
  },
  ioOutputGroups
};

export default genReducers(reducers, initState);
