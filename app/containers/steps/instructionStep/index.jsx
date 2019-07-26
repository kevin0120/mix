import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/es/Button/Button';
import { instructionStepActions } from '../../../modules/step/instructionStep/action';
import { StepContent } from '../types';
import { stepPayload, processingStep } from '../../../modules/order/selector';
import { Typography } from '@material-ui/core';

const mapState = (state, props) => ({
  ...props
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
