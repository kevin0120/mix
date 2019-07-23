// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles, Paper } from '@material-ui/core';
import stepTypes from '../steps/stepTypes';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';

const mapState = (state, props) => ({
  ...props,
  step: orderSelectors.viewingStep(state.order) || {},
  processingStep: orderSelectors.processingStep(state.order) || {},
  result: orderSelectors.stepData(orderSelectors.processingStep(state.order))?.result,
});
const mapDispatch = {};


const StepPageContainer = ({ step, processingStep, bindAction,result }) => {
  let stepProps = {};
  const classes = makeStyles(styles.stepPageContainer)();
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;
    stepProps =
      stepTypes?.[step?.type]?.genProps?.({ payload: step.payload || {} }) || stepProps;

    return (
      <div className={classes.root}>
        <Paper square className={classes.left}>

          <Paper square className={classes.image}>
            {
              (StepComponent && (
                <StepComponent
                  step={step}
                  {...stepProps}
                  isCurrent={step === processingStep}
                  bindAction={bindAction}
                />
              )) ||
              null
            }
          </Paper>
        </Paper>
        <Paper square className={classes.right}>
          <Paper square className={classes.description}>
            {step.description}
          </Paper>
          <Paper square className={classes.result}>
            {JSON.stringify(result)}
          </Paper>
        </Paper>
      </div>
    );
  }
  return null;
};

export default connect(mapState, mapDispatch)(StepPageContainer);
