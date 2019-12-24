import {STEP_ACTIONS} from './constants';

export default {
  submit:(result)=>({
    type:STEP_ACTIONS.SUBMIT,
    result
  })
}
