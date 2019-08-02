// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper } from '@material-ui/core';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';
import type { tOrder } from '../../modules/order/model';

type Props = {
  viewingOrder: tOrder
};

function StepWorking({ viewingOrder }: Props) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
    <div className={classes.root}>
      <Paper square className={classes.orderInfoContainer}>
        <Typography variant="h5">{viewingOrder?.name || ''}</Typography>
      </Paper>
    <div className={classes.main}>
      <Paper
        square
        classes={{ root: classes.leftContainer }}
      >
        <ButtonsContainer action={action}/>
        <StepPageContainer bindAction={bindAction}/>
      </Paper>
      <div className={classes.rightContainer}>
        <Paper square className={classes.stepperContainer}>
          <StepperContainer/>
        </Paper>
      </div>


    </div>
    </div>
  );
}

const mapState = (state, props) => ({
  ...props,
  viewingOrder: orderSelectors.viewingOrder(state.order) || {}
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
