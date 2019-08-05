export const SCREW_STEP = {
  RESULT: 'SCREW_STEP_RESULT'
};

export default {
  result: result => ({
    type: SCREW_STEP.RESULT,
    result
  })
};
