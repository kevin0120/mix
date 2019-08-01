// @flow
import {  makeStyles, StepContent, Typography } from '@material-ui/core';
import { connect } from 'react-redux';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import React from 'react';
import STEP_STATUS from '../../modules/step/model';
import styles from './styles';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { Dispatch } from '../../modules/indexReducer';
import type { tStepArray } from '../../modules/order/model';

const mapState = (state, props) => ({
  ...props,
  steps: orderSelectors.orderSteps(orderSelectors.viewingOrder(state.order)) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0
});

const mapDispatch = {
  jumpTo: orderActions.jumpToStep
};

type StepperLayoutProps = {
  steps: tStepArray,
  viewingIndex: {},
  jumpTo: Dispatch
};

const StepperContainer = ({
                            steps,
                            viewingIndex,
                            jumpTo
                          }: StepperLayoutProps) => {
  const classes = makeStyles(styles.stepperContainer)();
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
            error: s.status === STEP_STATUS.FAIL
          };
          return (
            <Step key={s.name} {...stepProps}>
              <StepButton
                completed={s.status === STEP_STATUS.FINISHED}
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

export default connect(
  mapState,
  mapDispatch
)(StepperContainer);
