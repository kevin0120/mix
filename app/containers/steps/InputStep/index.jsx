import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/es/Button/Button';
import { inputStepActions } from '../../../modules/steps/inputStep/action';

type Props = {
  label: string,
  isCurrent: boolean,
  submit: ()=>{},
  bindAction: ()=>{},
  step: {}
};

function InputStep({ step, label, isCurrent, submit, bindAction }: Props) {
  const [value, setValue] = useState('');

  const onSubmit = (v) => {
    setValue('');
    submit(v);
  };

  useEffect(() => {
    bindAction(<Button
      type="button"
      onClick={() => onSubmit(value)}
      disabled={!isCurrent}
    >submit</Button>);
    return () => bindAction(null);
  }, [step]);

  useEffect(() => {
    setValue('');
  }, [step]);

  return <div>
    {label}
    <input
      onChange={(e) => {
        setValue(e.target.value);
      }}
      value={value}
    />
  </div>;
}


const mapState = (state, props) => ({
  ...props
});

const mapDispatch = {
  submit: inputStepActions.submit
};

export default connect(mapState, mapDispatch)(InputStep);
