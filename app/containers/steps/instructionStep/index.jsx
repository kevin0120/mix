import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '../../../components/CustomButtons/Button';
import { Typography } from '@material-ui/core';
import { instructionStepActions } from '../../../modules/step/instructionStep/action';
import { StepContent } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  instruction: stepPayload(viewingStep(state.order))?.instruction || null
});

const mapDispatch = {
  submit: instructionStepActions.submit
};

type Props = {
  submit: () => {}
};

function InstructionStep({ step, isCurrent, submit, bindAction, instruction }: Props | StepContent) {

  useEffect(() => {
    bindAction(
      <Button
        type="button"
        color="primary"
        onClick={() => {
          submit();
        }}
        disabled={!isCurrent}
      >
        submit
      </Button>
    );
    return () => bindAction(null);
  }, [step, bindAction, isCurrent, submit]);

  return (
    <Typography>
      {instruction}
    </Typography>
  );
}

export default connect(mapState, mapDispatch)(InstructionStep);
