import {STEP_ACTIONS} from './constants';

export default {
  submit:()=>({
    type:STEP_ACTIONS.SUBMIT,
  }),
  input:(input)=>({
    type:STEP_ACTIONS.INPUT,
    input
  })
}
