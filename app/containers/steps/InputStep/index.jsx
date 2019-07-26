import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/es/Button/Button';
import { inputStepActions } from '../../../modules/step/inputStep/action';
import { StepContent } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  label: stepPayload(viewingStep(state.order))?.label || ''
});

const mapDispatch = {
  submit: inputStepActions.submit
};


type Props = {
  label: string,
  submit: () => {}
};

function InputStep({
                     step,
                     label,
                     isCurrent,
                     submit,
                     bindAction
                   }: Props | StepContent) {
  const [value, setValue] = useState('');

  useEffect(
    () => {
      const onSubmit = v => {
        submit(v);
      };
      bindAction(
        <Button
          type="button"
          onClick={() => {
            onSubmit(value);
          }}
          disabled={!isCurrent}
        >
          submit
        </Button>
      );
      return () => bindAction(null);
    },
    [step, bindAction, isCurrent, value, submit]
  );

  useEffect(
    () => {
      setValue('');
    },
    [step]
  );

  return (
    <div>
      {label}
      <input
        onChange={e => {
          setValue(e.target.value);
        }}
        value={value}
      />
    </div>
  );
}


export default connect(
  mapState,
  mapDispatch
)(InputStep);
