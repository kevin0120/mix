// @flow

export const TOOLS = {
  ENABLE: 'TOOLS.ENABLE',
  DISABLE: 'TOOLS.DISABLE',
  STATUS_CHANGE: 'TOOLS.STATUS.CHANGE',
  NEW_RESULTS: 'TOOLS_NEW_RESULTS'
};

export function toolEnable(reason) {
  return {
    type: TOOLS.ENABLE,
    enable: true,
    reason
  };
}

export function toolDisable(reason) {
  return {
    type: TOOLS.ENABLE,
    enable: false,
    reason
  };
}

export function toolStatusChange(toolSN, status, reason) {
  return {
    type: TOOLS.STATUS_CHANGE,
    toolSN,
    status,
    reason
  };
}

export function toolNewResults(results) {
  return {
    type: TOOLS.NEW_RESULTS,
    results,
  };
}
