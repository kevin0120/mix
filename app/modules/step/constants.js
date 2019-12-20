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
  scanner: 'scanner',
  instruction: 'text',
  text: 'text',
  screw: 'tightening',
  material: 'register_consumed_materials',
  check: 'pass_fail',
  video: 'video',
});

export default STEP_STATUS;
