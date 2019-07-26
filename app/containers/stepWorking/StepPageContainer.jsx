// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles, Paper, Grid, Typography } from '@material-ui/core';
import stepTypes from '../steps/stepTypes';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import type { tStep } from '../../modules/order/model';

const mapState = (state, props) => ({
  ...props,
  step: orderSelectors.viewingStep(state.order) || {},
  processingStep: orderSelectors.processingStep(state.order) || {},
  result: orderSelectors.stepData(orderSelectors.processingStep(state.order))?.result
});
const mapDispatch = {};

type Props = {
  step: tStep,
  processingStep: tStep,
  bindAction: ()=>{},
  result: any
};

const StepPageContainer = ({ step, processingStep, bindAction, result }): Props => {
  const classes = makeStyles(styles.stepPageContainer)();
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;

    return (
      <Grid container spacing={1} className={classes.root}>
        <Grid item container className={classes.left} spacing={1}>
          <Grid item className={classes.left}>
            <Paper square className={classes.image}>
              {
                (StepComponent && (
                  <StepComponent
                    step={step}
                    isCurrent={step === processingStep}
                    bindAction={bindAction}
                  />
                )) || null
              }
            </Paper>
          </Grid>
        </Grid>
        <Grid item container spacing={1} className={classes.right} direction={'column'}>

          <Grid item className={classes.description}>
            <Paper square className={classes.Paper}>
              <Typography>
                {step.description}
              </Typography>
            </Paper>
          </Grid>
          <Grid item className={classes.result}>
            <Paper square className={classes.Paper}>
              {JSON.stringify(result)}
            </Paper>
          </Grid>
        </Grid>
      </Grid>
    );
  }
  return null;
};

export default connect(mapState, mapDispatch)(StepPageContainer);
