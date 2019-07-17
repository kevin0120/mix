import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/es/Button/Button';
import { inputStepActions } from '../../../modules/steps/inputStep/action';
import { StepContent } from '../types';

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

const mapState = (state, props) => ({
  ...props
});

const mapDispatch = {
  submit: inputStepActions.submit
};

export default connect(
  mapState,
  mapDispatch
)(InputStep);
