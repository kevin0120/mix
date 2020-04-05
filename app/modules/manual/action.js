// @flow
import type { Result } from './types';

export const TIMELINE_STORY = {
  NEW: 'TIMELINE_STORY_NEW',
  CLEAR: 'TIMELINE_STORY_CLEAR'
};

export const MANUAL = {
  START: 'MANUAL_START',
  CLOSE: 'MANUAL_CLOSE',
  CLEAR: 'MANUAL_CLEAR',
  TIGHTENING: 'MANUAL_TIGHTENING',
  NO_SCANNER:'NO_SCANNER',
  SCANNER_NEW_DATA:'SCANNER_NEW_DATA',
  SELECT_TOOL: 'MANUAL_SELECT_TOOL',
  SELECT_PSET: 'SELECT_PSET',
  SET_DATA:'SET_DATA',
  GET_RESULT:'GET_RESULT',
  DID_MOUNT: 'DID_MOUNT',
  RESULT_INPUT:'RESULT_INPUT',
  INPUT_OK:'INPUT_OK'
};

export function inputOk1() {
  return {
    type: MANUAL.INPUT_OK
  }
}

export function resultInput(resultIn: Result) {
  return {
    type: MANUAL.RESULT_INPUT,
    resultIn
  };
}

export function clearStories() {
  return {
    type: TIMELINE_STORY.CLEAR
  };
}

export function didMount() {
  return {
    type: MANUAL.DID_MOUNT
  };
}


export function getresult(result) {
  return {
    type: MANUAL.GET_RESULT,
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

export function manualClear() {
  return {
    type: MANUAL.CLEAR
  };
}

export function noScanner() {
  return {
    type: MANUAL.NO_SCANNER
  };
}

export function newData(scanner) {
  return {
    type: MANUAL.SCANNER_NEW_DATA,
    scanner
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
    type: MANUAL.SET_DATA,
    ControllerSN,
    tool,
    pset
  };
}

