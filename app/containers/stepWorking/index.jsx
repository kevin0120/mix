import React from 'react';
import { connect } from 'react-redux';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import {
  orderActions
} from '../../modules/order/action';
import InputStep from '../../components/InputStep';

const stepTypes = {
  input: {
    component: InputStep,
    props: (props) => ({
      onSubmit: (value) => props.pushState(value, props.parallelId),
      label: props.payload.label
    })
  },
};

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

  renderStepContents = (step, isCurrent, parallelId) => {
    const { work } = this.props;
    let StepContent = null;
    let props = {};
    if (step && step.type && stepTypes[step.type] && stepTypes[step.type].component) {
      StepContent = stepTypes[step.type].component;
      props = stepTypes[step.type].props && stepTypes[step.type].props({
        parallelId,
        pushState: work,
        payload: step.payload || {}
      }) || props;
      return StepContent && <StepContent {...props} isCurrent={isCurrent} stepStatus={step.status || 'ready'}/>;
    }
    if (step instanceof Array) {
      return step.map((s, id) => this.renderStepContents(s, isCurrent, id));
    }
  };

  render() {
    const {
      order,
      trigger,
      next,
      previous
    } = this.props;
    const { currentOrder, currentViewingIndex, currentViewingStep, currentProcessingStep } = order;
    return <div>
      {(currentOrder &&
        currentOrder.steps &&
        this.renderSteps(currentOrder.steps, currentViewingIndex)) || null}
      <button type="button" onClick={() => previous()}>view previous</button>
      <button type="button" onClick={() => next()}>view next</button>
      <div>
        {JSON.stringify(currentViewingStep)}
        {
          this.renderStepContents(currentViewingStep, currentProcessingStep === currentViewingStep)
        }
      </div>
    </div>;
  }
}

const mapState = (state, props) => ({
  order: state.order,
  ...props
});

const mapDispatch = {
  trigger: orderActions.trigger,
  next: orderActions.nextStep,
  previous: orderActions.previousStep,
  work: orderActions.tryPushStep,
  jumpTo: orderActions.jumpToStep
};

export default connect(mapState, mapDispatch)(ConnectedStepWorking);
