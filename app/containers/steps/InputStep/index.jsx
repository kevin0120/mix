// @flow

import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '../../../components/CustomButtons/Button';
import { inputStepActions } from '../../../modules/step/inputStep/action';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import type { tStepProps } from '../types';

const mapState = (state, props) => ({
  ...props,
  label: stepPayload(viewingStep(state.order))?.label || ''
});

const mapDispatch = {
  submit: inputStepActions.submit
};


type Props = {
  label: string,
  submit: (?string | ?number) => any
} ;

function InputStep({
                     step,
                     label,
                     isCurrent,
                     submit,
                     bindAction
                   }: Props & tStepProps) {
  const [value, setValue] = useState('');

  useEffect(
    () => {
      const onSubmit = v => {
        submit(v);
      };
      bindAction(
        <Button
          type="button"
          color="primary"
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
