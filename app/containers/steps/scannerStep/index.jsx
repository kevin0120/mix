import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/es/Button/Button';
import { SvgIcon } from '@material-ui/core';
import styles from './styles';
import { scannerStepAction } from '../../../modules/steps/scannerStep/action';
import QRCode from './qrcode-scan.svg';

type Props = {
  classes: {},
  submit: ()=>{}
};

class scannerStep extends React.Component<Props> {


// eslint-disable-next-line flowtype/no-weak-types
  render(): React.ReactElement<any> | string | number | {} | React.ReactNodeArray | React.ReactPortal | boolean | null | undefined {
    const { classes, submit, isCurrent } = this.props;
    return <div className={classes.root}>
      <QRCode width="200" height="200" viewBox="0 0 24 24"/>
      <Button onClick={() => submit()} disabled={!isCurrent}>submit</Button>
    </div>;
  }
}

const mapState = (state, props) => ({
  ...props

});

const mapDispatch = {
  submit: scannerStepAction.submit
};

export default withStyles(styles)(connect(mapState, mapDispatch)(scannerStep));
