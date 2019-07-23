export const SCANNER_STEP = {
  GET_VALUE: 'SCANNER_STEP_GET_VALUE',
  SUBMIT: 'SCANNER_STEP_SUBMIT'
};

const getValue = (value) => ({
  type: SCANNER_STEP.GET_VALUE,
  value
});

const submit = () => ({
  type: SCANNER_STEP.SUBMIT
});

export const scannerStepAction = {
  getValue,
  submit
};
