// @flow

export const STEP_STATUS = Object.freeze({
  READY: 'ready',
  ENTERING: 'entering',
  DOING: 'doing',
  LEAVING: 'leaving',
  FAIL: 'fail',
  FINISHED: 'finished'
});

export const stepTypeKeys=Object.freeze({
  input: 'input',
  scanner: 'register_byproducts',
  instruction: 'instruction',
  text: 'text',
  screw: 'tightening',
  material: 'register_consumed_materials',
  passFail: 'passfail',
  measure: 'measure',
  video: 'video',
});

export const STEP_ACTIONS={
  SUBMIT:'STEP_SUBMIT',
  INPUT:'STEP_INPUT'
};

export default STEP_STATUS;
