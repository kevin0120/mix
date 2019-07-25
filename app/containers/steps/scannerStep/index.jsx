import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import makeStyles from '@material-ui/core/styles/makeStyles';
import styles from './styles';
import { scannerStepAction } from '../../../modules/step/scannerStep/action';
import QRCode from './qrcode-scan.svg';
import { StepContent } from '../types';
import withKeyboard from '../../../components/Keyboard';
import type { Dispatch } from '../../../modules/indexReducer';
import { stepData, processingStep } from '../../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  result: stepData(processingStep(state.order))?.result || {}
});

const mapDispatch = {
  submit: scannerStepAction.submit,
  getValue: scannerStepAction.getValue
};

type Props = {
  submit: Dispatch,
  getValue: Dispatch,
  keyboardInput: ()=>{}
};

function ScannerStep(
  {
    step,
    submit,
    isCurrent,
    bindAction,
    keyboardInput,
    label,
    getValue,
    result
  }: Props | StepContent) {
  const classes = makeStyles(styles)();
  useEffect(
    () => {
      bindAction(
        <Button onClick={() => submit()} disabled={!isCurrent}>
          submit
        </Button>
      );
      return () => bindAction(null);
    },
    [bindAction, isCurrent, step, submit]
  );

  return (
    <div className={classes.root}>
      <QRCode width="200" height="200" viewBox="0 0 24 24" fill={'#444'} />
      <TextField
        label={label}
        margin="normal"
        value={result?.[label] || ''}
        onClick={() => {
          keyboardInput({
            onSubmit: text => {
              getValue(text);
            },
            text: result?.[label] || '',
            title: label,
            label
          });
        }}
      />
    </div>
  );
}

export default connect(
  mapState,
  mapDispatch
)(withKeyboard(ScannerStep));
