import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import { StylesProvider } from '@material-ui/styles';

import Button from '@material-ui/core/es/Button/Button';
import { Typography, Paper, StepContent } from '@material-ui/core';
import { orderActions } from '../../modules/order/action';
import * as orderSelectors from '../../modules/order/selector';
import stepTypes from '../steps/stepTypes';
import styles from './styles';
import { ORDER_STEP_STATUS } from '../../modules/order/model';
import Dialog from '../../components/Dialog';
import type { Dispatch } from '../../modules/indexReducer';

const StepPage = ({ step, isCurrent, bindAction }) => {
  let stepProps = {};
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;
    stepProps =
      stepTypes?.[step?.type]?.genProps?.({
        payload: step.payload || {}
      }) || stepProps;
    return (
      (StepComponent && (
        <StepComponent
          step={step}
          {...stepProps}
          isCurrent={isCurrent}
          stepStatus={step.status || 'ready'}
          bindAction={bindAction}
        />
      )) ||
      null
    );
  }
  return null;
};

type StepperLayoutProps = {
  steps: [],
  currentStep: {},
  jumpTo: () => {}
};

const StepperLayout = ({ steps, currentStep, jumpTo }: StepperLayoutProps) => {
  const classes = makeStyles(styles.stepper)();
  return (
    <Stepper
      nonLinear
      activeStep={currentStep}
      orientation="vertical"
      className={classes.stepper}
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

const renderTimer = () => null;

type StepWorkingProps = {
  currentOrder: {},
  viewingStep: {},
  processingStep: {},
  viewingIndex: number,
  steps: [],
  next: Dispatch,
  previous: Dispatch,
  jumpTo: Dispatch,
  doNextStep: Dispatch,
  doPreviousStep: Dispatch
};

function StepWorking({
  currentOrder,
  viewingStep,
  processingStep,
  viewingIndex,
  steps,
  next,
  previous,
  jumpTo,
  doNextStep,
  doPreviousStep
}: StepWorkingProps) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;

  return (
    <StylesProvider injectFirst>
      <div className={classes.root}>
        <Dialog />
        <Paper
          square
          className={classes.leftContainer}
          classes={{ root: classes.leftContainer }}
        >
          <Paper square className={classes.orderInfoContainer}>
            <Typography variant="h5">{currentOrder?.name}</Typography>
          </Paper>
          <div className={classes.buttonsContainer}>
            <div>
              <Button
                disabled={noPrevious}
                type="button"
                onClick={() => previous()}
              >
                {'<<'}
              </Button>
              <Button disabled={noNext} type="button" onClick={() => next()}>
                {'>>'}
              </Button>
              {viewingStep?.skippable && (
                <Button
                  disabled={noNext || viewingStep !== processingStep}
                  type="button"
                  onClick={() => doNextStep()}
                >
                  {'skip'}
                </Button>
              )}
              {viewingStep?.undoable && (
                <Button
                  disabled={viewingStep !== processingStep}
                  type="button"
                  onClick={() => doPreviousStep()}
                >
                  {'undo'}
                </Button>
              )}
            </div>
            <div>{action}</div>
          </div>
          <div className={classes.contentContainer}>
            <StepPage
              step={viewingStep}
              isCurrent={viewingStep === processingStep}
              bindAction={bindAction}
            />
          </div>
        </Paper>
        <div className={classes.rightContainer}>
          <Paper square className={classes.timerContainer}>
            {renderTimer()}
          </Paper>
          <Paper square className={classes.stepperContainer}>
            <StepperLayout
              steps={steps}
              currentStep={viewingIndex}
              jumpTo={jumpTo}
            />
          </Paper>
        </div>
      </div>
    </StylesProvider>
  );
}

const mapState = (state, props) => ({
  currentOrder: orderSelectors.currentOrder(state.order) || {},
  viewingStep: orderSelectors.viewingStep(state.order) || {},
  processingStep: orderSelectors.processingStep(state.order) || {},
  steps: orderSelectors.orderSteps(state.order) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0,
  ...props
});

const mapDispatch = {
  next: orderActions.nextStep,
  doNextStep: orderActions.doNextStep,
  previous: orderActions.previousStep,
  jumpTo: orderActions.jumpToStep,
  doPreviousStep: orderActions.doPreviousStep
};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
