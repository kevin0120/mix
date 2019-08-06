import type { tInput } from '../../../common/type';

export const SCANNER_STEP = {
  GET_VALUE: 'SCANNER_STEP_GET_VALUE',
  SUBMIT: 'SCANNER_STEP_SUBMIT'
};

const getValue = (input: tInput) => ({
  type: SCANNER_STEP.GET_VALUE,
  input
});

const submit = () => ({
  type: SCANNER_STEP.SUBMIT
});

export const scannerStepAction = {
  getValue,
  submit
};
