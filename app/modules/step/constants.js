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
  instruction: 'text',
  text: 'text',
  screw: 'tightening',
  material: 'register_consumed_materials',
  passFail: 'pass_fail',
  measure: 'measure',
  video: 'video',
});

export default STEP_STATUS;
