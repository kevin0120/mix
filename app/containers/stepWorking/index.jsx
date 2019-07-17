import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import Button from '@material-ui/core/es/Button/Button';
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
      />;
  }
};

const StepperLayout = ({ steps, currentStep, onClick }) => {
  const classes = makeStyles(styles.stepper)();
  return (
    <Stepper nonLinear activeStep={currentStep} orientation="vertical" className={classes.stepper}>
      {steps.map((s, id) => {
        const stepProps = {};
        const labelProps = {};
        return (
          <Step key={s.name} {...stepProps}>
            <StepButton completed={s.status === 'finish'} onClick={() => onClick(id)} className={classes.stepButton}>
              <StepLabel {...labelProps}>{s.name}</StepLabel>
            </StepButton>
          </Step>
        );
      })}
    </Stepper>
  );
};

const renderTimer = () => 'here is a timer';

function StepWorking({ order, next, previous, jumpTo }) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
    <div className={classes.root}>
      <div className={classes.leftContainer}>
        <div className={classes.orderInfoContainer}>
        </div>
        <div className={classes.buttonsContainer}>
          <div>
            <Button type="button" onClick={() => previous()}>{'<<'}</Button>
            <Button type="button" onClick={() => next()}>{'>>'}</Button>
          </div>
          <div>
            {action}
          </div>
        </div>
        <div className={classes.contentContainer}>
          <StepContents
            step={viewingStep(order)}
            isCurrent={viewingStep(order) === processingStep(order)}
            bindAction={bindAction}
          />
        </div>
      </div>
      <div className={classes.rightContainer}>
        <div className={classes.timerContainer}>
          {renderTimer()}
        </div>
        <div className={classes.stepperContainer}>
          {(currentOrder(order) &&
            currentOrder(order).steps &&
            <StepperLayout
              steps={orderSteps(order)}
              currentStep={viewingIndex(order)}
              jumpTo={jumpTo}
            />
          ) || null}
        </div>
      </div>
    </div>
  );
}

const mapState = (state, props) => ({
  order: state.order,
  ...props
});

const mapDispatch = {
  next: orderActions.nextStep,
  previous: orderActions.previousStep,
  jumpTo: orderActions.jumpToStep
};

export default connect(mapState, mapDispatch)(StepWorking);
