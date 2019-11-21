import { IO } from './constants';

export default {
  set: (group, status) => ({
    type: IO.SET,
    group,
    status
  }),
  addListener: (inputType, action) => ({
    type: IO.ADD_LISTENER,
    inputType, action
  }),
  setPort: (output, portIdx) => ({
    type: IO.SET_PORT,
    output,
    portIdx
  })
};
