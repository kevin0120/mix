export const NETWORK = {
  CHECK: 'NETWORK_CHECK',
  CHECK_OK: 'NETWORK_CHECK_OK',
  SET: 'NETWORK_SET',
  SET_OK: 'NETWORK_SET_OK',
  SET_FAIL: 'NETWORK_SET_FAIL',
  SCAN: 'NETWORK_SCAN',
  SCAN_OK: 'NETWORK_SCAN_OK',
  ERROR: 'NETWORK_ERROR',
  SIGNAL:'NETWORK_SIGNAL',
  SIGNAL_OK:'NETWORK_SIGNAL_OK'
};

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
