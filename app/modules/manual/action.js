// @flow
export const MANUAL = {
  START: 'MANUAL_START',
  CLOSE: 'MANUAL_CLOSE',
  CANCEL: 'MANUAL_CANCEL',
  TIGHTENING: 'MANUAL_TIGHTENING',
  NOSCANNER:'NO_SCANNER',
  NEWCANNERDATA:'SCANNER_NEW_DATA',
  SELECT_TOOL: 'MANUAL_SELECT_TOOL',
  SELECT_PSET: 'MANUAL_SELECT_PSET',
  SETDATA:'SET_DATA',
  GETRESULT:'GET_RESULT',
  DIDMOUNT: 'DID_MOUNT',
};

export function didMount() {
  return {
    type: MANUAL.DIDMOUNT
  };
}


export function getresult(result) {
  return {
    type: MANUAL.GETRESULT,
    result
  };
}


export function start() {
  return {
    type: MANUAL.START
  };
}

export function manualTightening() {
  return {
    type: MANUAL.TIGHTENING
  };
}

export function close() {
  return {
    type: MANUAL.CLOSE
  };
}



export function noScanner() {
  return {
    type: MANUAL.NOSCANNER
  };
}

export function newData(scanner) {
  return {
    type: MANUAL.NEWCANNERDATA,
    scanner
  };
}

export function manualCancel() {
  return {
    type: MANUAL.CANCEL
  };
}
export function selectTool(tool){
  return {
    type: MANUAL.SELECT_TOOL,
    tool
  };
}


export function selectPset(pset) {
  return {
    type: MANUAL.SELECT_PSET,
    pset
  };
}

export function setData(ControllerSN: string,tool: string,pset: number) {
  return {
    type: MANUAL.SETDATA,
    ControllerSN,
    tool,
    pset
  };
}

