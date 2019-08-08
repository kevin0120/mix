export const MATERIAL_STEP = {
  IO_INPUT: 'MATERIAL_STEP_IO_INPUT'
};

export default {
  ioInput: (input) => ({
    type: MATERIAL_STEP.IO_INPUT,
    input
  })
};
