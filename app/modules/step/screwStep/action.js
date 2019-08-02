export const SCREW_STEP = {
  RESULT: 'SCREW_STEP_RESULT',
  IMAGE_READY: 'SCREW_STEP_IMAGE_READY'
};

export default {
  imageReady: () => ({
    type: SCREW_STEP.IMAGE_READY
  }),
  result: (result) => ({
    type: SCREW_STEP.RESULT,
    result
  })
};
