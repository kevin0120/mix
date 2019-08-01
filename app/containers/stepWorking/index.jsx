// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper, createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import { StylesProvider } from '@material-ui/styles';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import Dialog from '../../components/Dialog';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';
import Timer from './Timer';
import type { tOrder } from '../../modules/order/model';

type StepWorkingProps = {
  viewingOrder: tOrder
};

function StepWorking({ viewingOrder }: StepWorkingProps) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
        <div className={classes.root}>
          <Dialog/>
          <Paper
            square
            className={classes.leftContainer}
            classes={{ root: classes.leftContainer }}
          >
            <Paper square className={classes.orderInfoContainer}>
              <Typography variant="h5">{viewingOrder?.name}</Typography>
            </Paper>
            <ButtonsContainer action={action}/>
            <StepPageContainer bindAction={bindAction}/>
          </Paper>
          <div className={classes.rightContainer}>
            <Paper square className={classes.timerContainer}>
              <Timer/>
            </Paper>
            <Paper square className={classes.stepperContainer}>
              <StepperContainer/>
            </Paper>
          </div>
        </div>
  );
}

const mapState = (state, props) => ({
  ...props,
  viewingOrder: orderSelectors.viewingOrder(state.order) || {},
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
