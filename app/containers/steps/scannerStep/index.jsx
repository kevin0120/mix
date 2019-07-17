import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import Button from '@material-ui/core/es/Button/Button';
import makeStyles from '@material-ui/core/styles/makeStyles';
import styles from './styles';
import { scannerStepAction } from '../../../modules/steps/scannerStep/action';
import QRCode from './qrcode-scan.svg';

type Props = {
  submit: ()=>{},
  isCurrent: boolean,
  step: {},
  bindAction: ()=>{}
};

function ScannerStep({ step, submit, isCurrent, bindAction }: Props) {
  const classes = makeStyles(styles)();

  useEffect(() => {
    bindAction(<Button onClick={() => submit()} disabled={!isCurrent}>submit</Button>);
    return () => bindAction(null);
  }, [bindAction, isCurrent, step, submit]);

  return <div className={classes.root}>
    <QRCode width="200" height="200" viewBox="0 0 24 24"/>
  </div>;
}

const mapState = (state, props) => ({
  ...props
});

const mapDispatch = {
  submit: scannerStepAction.submit
};

export default connect(mapState, mapDispatch)(ScannerStep);
