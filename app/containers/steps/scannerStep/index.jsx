// @flow

import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import TextField from '@material-ui/core/TextField';
import makeStyles from '@material-ui/core/styles/makeStyles';
// import QRcode from 'qrcode.react';
import Button from '../../../components/CustomButtons/Button';
import styles from './styles';
import { scannerStepAction } from '../../../modules/step/scannerStep/action';
import type { tStepProps } from '../types';
import withKeyboard from '../../../components/Keyboard';
import type { Dispatch } from '../../../modules/typeDef';
import { stepData, stepPayload, viewingStep } from '../../../modules/order/selector';
import type { tKeyboardBinder } from '../../../components/Keyboard/typeDef';
import ScanIcon from './barcode-scan.svg';

type tDP = {|
  submit: Dispatch,
  getValue: Dispatch
|};
type tSP = {|
  ...tOP,
  result: string,
  label: string
|};

type tOP = {|
  ...tStepProps,
  keyboardInput: tKeyboardBinder
|};

type Props = {
  ...tOP,
  ...tDP,
  ...tSP
};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  result: stepData(viewingStep(state.order))?.result || '',
  label: stepPayload(viewingStep(state.order))?.label || ''
});

const mapDispatch: tDP = {
  submit: scannerStepAction.submit,
  getValue: scannerStepAction.getValue
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
                     }: Props) {
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
      {/*<QRcode value="This Is Demo QR Code" size={400}/>*/}
      <ScanIcon width="200" height="200" viewBox="0 0 24 24" fill="#424242"/>
      <TextField
        className={classes.input}
        label={label || '请扫描或输入条码'}
        disabled={!isCurrent}
        margin="normal"
        value={result || ''}
        variant="outlined"
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
              text: result || '',
              title: label,
              label
            });
          }
        }}
      />
    </div>
  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(withKeyboard(ScannerStep));
