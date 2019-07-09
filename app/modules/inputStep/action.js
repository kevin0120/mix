export const INPUT_STEP={
  SUBMIT:'INPUT_STEP_SUBMIT'
};

export const inputStepAtions={
  submit:(payload)=>({
    type:INPUT_STEP.SUBMIT,
    payload
  })
};
