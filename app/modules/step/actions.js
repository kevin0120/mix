import {STEP_ACTIONS} from './constants';

export default {
  submit:(result)=>({
    type:STEP_ACTIONS.SUBMIT,
    result
  }),
  input:(input)=>({
    type:STEP_ACTIONS.INPUT,
    input
  }),
  confirmFail:()=>({
    type:STEP_ACTIONS.CONFIRM_FAIL
  })
}
