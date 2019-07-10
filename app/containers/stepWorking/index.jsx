import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import Button from '@material-ui/core/es/Button/Button';
import {
  orderActions
} from '../../modules/order/action';
import stepTypes from './stepTypes';
import styles from './styles';


class ConnectedStepWorking extends React.Component {
  renderSteps(steps, currentStep) {
    const { jumpTo } = this.props;
    return <Stepper nonLinear activeStep={currentStep}>
      {steps.map((s, id) => {
        const stepProps = {};
        const labelProps = {};
        return (
          <Step key={s.name} {...stepProps}>
            <StepButton completed={s.status === 'finish'} onClick={() => jumpTo(id)}>
              <StepLabel {...labelProps}>{s.name}</StepLabel>
            </StepButton>
          </Step>
        );
      })}
    </Stepper>;
  }

  renderStepContents = (step, isCurrent) => {
    let StepContent = null;
    let props = {};
    if (step && step.type && stepTypes[step.type] && stepTypes[step.type].component) {
      StepContent = stepTypes[step.type].component;
      props = stepTypes[step.type].props && stepTypes[step.type].props({
        payload: step.payload || {}
      }) || props;
      return StepContent && <StepContent {...props} isCurrent={isCurrent} stepStatus={step.status || 'ready'}/>;
    }
  };

  render() {
    const {
      classes,
      order,
      next,
      previous
    } = this.props;
    const { currentOrder, currentViewingIndex, currentViewingStep, currentProcessingStep } = order;
    return <div className={classes.root}>
      {(currentOrder &&
        currentOrder.steps &&
        this.renderSteps(currentOrder.steps, currentViewingIndex)) || null}
      <Button type="button" onClick={() => previous()}>view previous</Button>
      <Button type="button" onClick={() => next()}>view next</Button>
      {/* {JSON.stringify(currentViewingStep)} */}
      <div className={classes.contentContainer}>
        {this.renderStepContents(currentViewingStep, currentProcessingStep === currentViewingStep)}
      </div>
    </div>;
  }
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

export default withStyles(styles)(connect(mapState, mapDispatch)(ConnectedStepWorking));
