// @flow
import { MANUAL, TIMELINE_STORY } from './action';
import type { StateType, Action } from './types';
import type { tAction, tWorkCenterMode } from '../workCenterMode/interface/typeDef';

const configs = {
  counter: 2,
  scanner: '',
  controllerSN: '',
  tool: '',
  pset: 0,
  result: [],
  timeline: [],
  workingTools:[]
};

const manualReducers = {
  [MANUAL.SET_DATA]: (state, action) => ({
    ...state,
    controllerSN: action.ControllerSN,
    tool: action.tool,
    pset: action.pset
  }),
  [MANUAL.SCANNER_NEW_DATA]: (state, action) => ({
    ...state,
    scanner: action.scanner
  }),
  [MANUAL.GET_RESULT]: (state, action) => ({
    ...state,
    result: state.result.concat(action.result)
  }),
  [MANUAL.START]: (state, action) => ({
    ...state,
    result: [],
    working: true
  }),
  [TIMELINE_STORY.NEW]: (state, action) => ({
    ...state,
    timeline: NewStory(state.timeline, action.story)
  }),
  [TIMELINE_STORY.CLEAR]: (state, action) => ({
    ...state,
    timeline: ClearStory(),
    working: false
  }),
  [MANUAL.RESULT_INPUT]: (state, action) => ({
    ...state,
    resultIn: action.resultIn
  })
};

export function NewStory(state, story) {
  return [story, ...state];
}

export function ClearStory() {
  return [];
}

export default function(state: StateType = configs, action: Action): tWorkCenterMode {
  if (manualReducers[(action?.type)]) {
    return manualReducers[(action?.type)](state, action);
  }
  return state;
}


