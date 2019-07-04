// @flow

export const TOOLS = {
  ENABLE: 'TOOLS.ENABLE',
  DISABLE: 'TOOLS.DISABLE',
  STATUS_CHANGE:'TOOLS.STATUS.CHANGE'
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

export function toolStatusChange(toolSN,status,reason){
  return{
    type:TOOLS.STATUS_CHANGE,
    toolSN,
    status,
    reason
  }
}
