import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import Button from '@material-ui/core/es/Button/Button';
import { Typography, Paper } from '@material-ui/core';
import { orderActions } from '../../modules/order/action';
import { currentOrder, orderSteps, processingStep, viewingIndex, viewingStep } from '../../modules/order/selector';
import stepTypes from '../steps/stepTypes';
import styles from './styles';


const StepContents = ({ step, isCurrent, bindAction }) => {
  let stepProps = {};
  if (step && step.type && stepTypes[step.type] && stepTypes[step.type].component) {
    const StepContent = stepTypes[step.type].component;
    stepProps = stepTypes[step.type].props && stepTypes[step.type].props({
      payload: step.payload || {}
    }) || stepProps;
    return StepContent &&
      <StepContent
        step={step}
        {...stepProps}
        isCurrent={isCurrent}
        stepStatus={step.status || 'ready'}
        bindAction={bindAction}
      />
      || null;
  }
  return null;
};

const StepperLayout = ({ steps, currentStep, jumpTo }) => {
  const classes = makeStyles(styles.stepper)();
  return (
    <Stepper nonLinear activeStep={currentStep} orientation="vertical" className={classes.stepper}>
      {steps.map((s, id) => {
        const stepProps = {};
        const labelProps = {};
        return (
          <Step key={s.name} {...stepProps}>
            <StepButton completed={s.status === 'finish'} onClick={() => jumpTo(id)} className={classes.stepButton}>
              <StepLabel {...labelProps}>{s.name}</StepLabel>
            </StepButton>
          </Step>
        );
      })}
    </Stepper>
  );
};

const renderTimer = () => 'here is a timer';

function StepWorking({ currentOrder, viewingStep, processingStep, viewingIndex, steps, next, previous, jumpTo }) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;

  return (
    <div className={classes.root}>
      <Paper square className={classes.leftContainer}>
        <div className={classes.orderInfoContainer}>
          <Typography variant="h5">
            {currentOrder && currentOrder.name}
          </Typography>
        </div>
        <div className={classes.buttonsContainer}>
          <div>
            <Button disabled={noPrevious} type="button" onClick={() => previous()}>{'<<'}</Button>
            <Button disabled={noNext} type="button" onClick={() => next()}>{'>>'}</Button>
          </div>
          <div>
            {action}
          </div>
        </div>
        <div className={classes.contentContainer}>
          <StepContents
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
  );
}

const mapState = (state, props) => {
  if (currentOrder(state.order)) {
    return {
      currentOrder: currentOrder(state.order),
      viewingStep: viewingStep(state.order),
      processingStep: processingStep(state.order),
      steps: orderSteps(state.order),
      viewingIndex: viewingIndex(state.order),
      ...props
    };
  }
  return {
    currentOrder: {},
    viewingStep: {},
    processingStep: {},
    steps: [],
    viewingIndex: 0,
    ...props
  };
};

const mapDispatch = {
  next: orderActions.nextStep,
  previous: orderActions.previousStep,
  jumpTo: orderActions.jumpToStep
};

export default connect(mapState, mapDispatch)(StepWorking);
