import React from 'react';
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


const renderStepContents = (step, isCurrent) => {
  let StepContent = null;
  let stepProps = {};
  if (step && step.type && stepTypes[step.type] && stepTypes[step.type].component) {
    StepContent = stepTypes[step.type].component;
    stepProps = stepTypes[step.type].props && stepTypes[step.type].props({
      payload: step.payload || {}
    }) || stepProps;
    return StepContent && <StepContent {...stepProps} isCurrent={isCurrent} stepStatus={step.status || 'ready'}/>;
  }
};

const renderSteps = (steps, currentStep, onClick) => (
  <Stepper nonLinear activeStep={currentStep}>
    {steps.map((s, id) => {
      const stepProps = {};
      const labelProps = {};
      return (
        <Step key={s.name} {...stepProps}>
          <StepButton completed={s.status === 'finish'} onClick={() => onClick(id)}>
            <StepLabel {...labelProps}>{s.name}</StepLabel>
          </StepButton>
        </Step>
      );
    })}
  </Stepper>
);

function StepWorking(props) {
  const { order, next, previous, jumpTo } = props;
  const classes = makeStyles(styles)();

  return (
    <div className={classes.root}>
      {(currentOrder(order) &&
        currentOrder(order).steps &&
        renderSteps(orderSteps(order), viewingIndex(order), jumpTo)) || null}
      <div className={classes.buttonsContainer}>
        <Button type="button" onClick={() => previous()}>view previous</Button>
        <Button type="button" onClick={() => next()}>view next</Button>
      </div>
      <div className={classes.contentContainer}>
        {renderStepContents(viewingStep(order), viewingStep(order) === processingStep(order))}
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
