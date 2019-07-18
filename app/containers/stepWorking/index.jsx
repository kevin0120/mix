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
// import CoumpoundTimer from 'react-compound-timer';
import { orderActions } from '../../modules/order/action';
import { currentOrder, orderSteps, processingStep, viewingIndex, viewingStep } from '../../modules/order/selector';
import stepTypes from '../steps/stepTypes';
import styles from './styles';

const StepPage = ({ step, isCurrent, bindAction }) => {
  let stepProps = {};
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;
    stepProps = stepTypes?.[step?.type]?.props({
      payload: step.payload || {}
    }) || stepProps;
    return StepComponent &&
      <StepComponent
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
            <StepContent>
              <Typography>
                {s.info}
              </Typography>
            </StepContent>
          </Step>
        );
      })}
    </Stepper>
  );
};

const renderTimer = () =>null;
// const renderTimer = () =>
//   <Typography variant="h2">
//     <CoumpoundTimer formatValue={(v) => `0${v}`.slice(-2)}>
//       <CoumpoundTimer.Hours/>:
//       <CoumpoundTimer.Minutes/>:
//       <CoumpoundTimer.Seconds/>
//     </CoumpoundTimer>
//   </Typography>;

function StepWorking({ currentOrder, viewingStep, processingStep, viewingIndex, steps, next, previous, jumpTo }) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;

  return (
    <StylesProvider injectFirst>
      <div className={classes.root}>
        <Paper square className={classes.leftContainer} classes={{ root: classes.leftContainer }}>
          <Paper square className={classes.orderInfoContainer}>
            <Typography variant="h5">
              {currentOrder?.name}
            </Typography>
          </Paper>
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
      currentOrder: currentOrder(state.order)||{},
      viewingStep: viewingStep(state.order)||{},
      processingStep: processingStep(state.order)||{},
      steps: orderSteps(state.order)||[],
      viewingIndex: viewingIndex(state.order)||0,
      ...props
    });

const mapDispatch = {
  next: orderActions.nextStep,
  previous: orderActions.previousStep,
  jumpTo: orderActions.jumpToStep
};

export default connect(mapState, mapDispatch)(StepWorking);
