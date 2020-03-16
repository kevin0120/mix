// @flow
import { MANUAL,TIMELINE_STORY} from './action';
import type { StateType, Action } from './types';

const configs = {
  counter: 2,
  scanner: '',
  controllerSN:'',
  tool: '',
  pset: 0,
  result:[],
  timeline:[]
}

export default function (state: StateType = configs, action: Action) {
  switch (action.type) {
    case MANUAL.SETDATA:
      return { ...state,controllerSN:action.ControllerSN,tool:action.tool, pset:action.pset };
    case MANUAL.NEWCANNERDATA:
      return { ...state, scanner: action.scanner };
    case MANUAL.GETRESULT:
      return { ...state, result: state.result.concat(action.result) };
    case MANUAL.START:
      return { ...state, working: true };
    case MANUAL.CLOSE:
      return { ...state, working: false};
    case MANUAL.DIDMOUNT:
      return { ...state, result: []};
    case TIMELINE_STORY.NEW:
      return { ...state, timeline: NewStory(state.timeline, action.story)};
    case TIMELINE_STORY.CLEAR:
      return { ...state, timeline: ClearStory()};
    case MANUAL.RESULTINPUT:
      return { ...state, resultIn: action.resultIn};
    default:
      return state;
  }
}


export function NewStory(state, story) {
  return [story, ...state];
}

export function ClearStory() {
  return [];
}
