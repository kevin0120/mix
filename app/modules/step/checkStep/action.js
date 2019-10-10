export const CHECK_STEP = {
  SUBMIT: 'CHECK_STEP_SUBMIT',
  CANCEL: 'CHECK_STEP_CANCEL'
};

export default {
  submit: payload => ({
    type: CHECK_STEP.SUBMIT,
    payload
  }),
  cancel: () => ({
    type: CHECK_STEP.CANCEL
  })
};
