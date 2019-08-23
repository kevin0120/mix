// @flow

import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import TextField from '@material-ui/core/TextField';
import makeStyles from '@material-ui/core/styles/makeStyles';
import QRcode from 'qrcode.react';
import Button from '../../../components/CustomButtons/Button';
import styles from './styles';
import { scannerStepAction } from '../../../modules/step/scannerStep/action';
// import QRCode from './qrcode-scan.svg';
import { StepContent } from '../types';
import withKeyboard from '../../../components/Keyboard';
import type { Dispatch } from '../../../modules/indexReducer';
import {
  stepData,
  stepPayload,
  viewingStep
} from '../../../modules/order/selector';
// import { scanner } from '../../../modules/scanner/saga';

const mapState = (state, props) => ({
  ...props,
  result: stepData(viewingStep(state.order))?.result || {},
  label: stepPayload(viewingStep(state.order))?.label || ''
});

const mapDispatch = {
  submit: scannerStepAction.submit,
  getValue: scannerStepAction.getValue
};

type Props = {
  submit: Dispatch,
  getValue: Dispatch,
  keyboardInput: Function
};

function ScannerStep({
  step,
  submit,
  isCurrent,
  bindAction,
  keyboardInput,
  label,
  getValue,
  result
}: Props & StepContent) {
  const classes = makeStyles(styles)();
  useEffect(() => {
    bindAction(
      <Button onClick={() => submit()} disabled={!isCurrent} color="primary">
        完成
      </Button>
    );
    return () => bindAction(null);
  }, [bindAction, isCurrent, step, submit]);

  return (
    <div className={classes.root}>
      <QRcode value="This Is Demo QR Code" size={400} />
      <TextField
        label={label}
        disabled={!isCurrent}
        margin="normal"
        value={result?.[label] || ''}
        onClick={() => {
          if (isCurrent) {
            keyboardInput({
              onSubmit: text => {
                getValue({
                  data: text,
                  time: new Date(),
                  name: 'input'
                });
              },
              text: result?.[label] || '',
              title: label,
              label
            });
          }
        }}
      />
    </div>
  );
}

export default connect(
  mapState,
  mapDispatch
)(withKeyboard(ScannerStep));
