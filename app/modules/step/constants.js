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
  instruction: 'instruction',
  screw: 'screw',
  material: 'material',
  check: 'check',
  video: 'video'
});

export default STEP_STATUS;
