import { makeStyles, StepContent, Typography } from '@material-ui/core';
import { useSelector, useDispatch } from 'react-redux';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import React from 'react';
import { ORDER_STEP_STATUS } from '../../modules/order/model';
import styles from './styles';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';

const StepperContainer = () => {
  const dispatch = useDispatch();
  const classes = makeStyles(styles.stepperContainer)();
  const steps = useSelector(state => orderSelectors.orderSteps(state.order) || []);
  const viewingIndex = useSelector(state => orderSelectors.viewingIndex(state.order) || 0);

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
              onClick={() => dispatch(orderActions.jumpToStep(id))}
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

export default StepperContainer;
