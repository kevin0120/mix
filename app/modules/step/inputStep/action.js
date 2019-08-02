export const INPUT_STEP = {
  SUBMIT: 'INPUT_STEP_SUBMIT'
};

export const inputStepActions = {
  submit: (payload) => ({
    type: INPUT_STEP.SUBMIT,
    payload
  })
};
