import { makeStyles, StepContent, Typography } from '@material-ui/core';
import { connect } from 'react-redux';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import React from 'react';
import { ORDER_STEP_STATUS } from '../../modules/order/model';
import styles from './styles';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';

const mapState = (state, props) => ({
  ...props,
  steps: orderSelectors.orderSteps(state.order) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0
});

const mapDispatch = {
  jumpTo: orderActions.jumpToStep
};

type StepperLayoutProps = {
  steps: [],
  viewingIndex: {},
  jumpTo: () => {}
};

const StepperContainer = ({ steps, viewingIndex, jumpTo }: StepperLayoutProps) => {
  const classes = makeStyles(styles.stepper)();
  return (
    <Stepper
      nonLinear
      activeStep={viewingIndex}
      orientation="vertical"
      className={classes.root}
    >
      {steps.map((s, id) => {
        const stepProps = {};
        const labelProps = {
          error: s.status === ORDER_STEP_STATUS.FAIL
        };
        return (
          <Step key={s.name} {...stepProps}>
            <StepButton
              completed={s.status === ORDER_STEP_STATUS.FINISHED}
              onClick={() => jumpTo(id)}
              className={classes.stepButton}
            >
              <StepLabel {...labelProps}>{s.name}</StepLabel>
            </StepButton>
            <StepContent>
              <Typography>{s.info}</Typography>
            </StepContent>
          </Step>
        );
      })}
    </Stepper>
  );
};

export default connect(mapState, mapDispatch)(StepperContainer);
