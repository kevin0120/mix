export const MATERIAL_STEP = {
  READY: 'MATERIAL_STEP_READY',
  ITEM: 'MATERIAL_STEP_ITEM'
};

export default {
  ready: input => ({
    type: MATERIAL_STEP.READY,
    input
  }),
  item: item => ({
    type: MATERIAL_STEP.ITEM,
    item
  })
};
