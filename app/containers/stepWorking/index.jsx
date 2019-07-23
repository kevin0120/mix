import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { StylesProvider } from '@material-ui/styles';
import { Typography, Paper } from '@material-ui/core';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import Dialog from '../../components/Dialog';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';


const renderTimer = () => null;

type StepWorkingProps = {
  currentOrder: {}
};

function StepWorking({ currentOrder }: StepWorkingProps) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
    <StylesProvider injectFirst>
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
            {renderTimer()}
          </Paper>
          <Paper square className={classes.stepperContainer}>
            <StepperContainer/>
          </Paper>
        </div>
      </div>
    </StylesProvider>
  );
}

const mapState = (state, props) => ({
  ...props,
  currentOrder: orderSelectors.currentOrder(state.order) || {}
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
