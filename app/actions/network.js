import { NETWORK } from './actionTypes';


export function networkScan() {
  return {
    type: NETWORK.SCAN
  };
}

export const networkCheck = () => ({
  type: NETWORK.CHECK
});

export function networkSet(data) {
  return {
    type: NETWORK.SET,
    data
  };
}

export function networkSignal() {
  return {
    type: NETWORK.SIGNAL,
  };
}
