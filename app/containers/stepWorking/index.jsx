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

type StepWorkingProps = {
  currentOrder: {}
};

const theme = createMuiTheme({
  overrides: {
    MuiPaper: {
      root: {
        backgroundColor: 'inherit'
      }
    }
  }
});

function StepWorking({ currentOrder }: StepWorkingProps) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
    <StylesProvider injectFirst>
      <MuiThemeProvider theme={theme}>
        <div className={classes.root}>
          <Dialog/>
          <Paper
            square
            className={classes.leftContainer}
            classes={{ root: classes.leftContainer }}
          >
            <Paper square className={classes.orderInfoContainer}>
              <Typography variant="h5">{currentOrder?.name}</Typography>
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
      </MuiThemeProvider>
    </StylesProvider>

  );
}

const mapState = (state, props) => ({
  ...props,
  currentOrder: orderSelectors.currentOrder(state.order) || {},
  startTime: orderSelectors.startTime(orderSelectors.viewingStep(state.order)) || null,
  endTime: orderSelectors.endTime(orderSelectors.viewingStep(state.order)) || null
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
