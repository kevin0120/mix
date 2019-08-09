export const MATERIAL_STEP = {
  READY: 'MATERIAL_STEP_READY'
};

export default {
  ready: (input) => ({
    type: MATERIAL_STEP.READY,
    input
  })
};
