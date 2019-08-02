// @flow
import { makeStyles, StepContent, Typography } from '@material-ui/core';
import { connect } from 'react-redux';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import React, { useState, useEffect } from 'react';
import ReactDOM from 'react-dom';
import { Loop } from '@material-ui/icons';
import STEP_STATUS from '../../modules/step/model';
import styles from './styles';
import * as oSel from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { Dispatch } from '../../modules/indexReducer';
import type { tStep, tStepArray } from '../../modules/order/model';
import Timer from './Timer';

const mapState = (state, props) => ({
  ...props,
  steps: oSel.orderSteps(oSel.viewingOrder(state.order)) || [],
  workingStep: oSel.workingStep(oSel.workingOrder(state.order)),
  viewingStep: oSel.viewingStep(state.order),
  viewingIndex: oSel.viewingIndex(state.order) || 0,
  isCurrent: oSel.viewingOrder(state.order) === oSel.workingOrder(state.order)
});

const mapDispatch = {
  jumpTo: orderActions.jumpToStep
};

type StepperLayoutProps = {
  steps: tStepArray,
  viewingIndex: number,
  jumpTo: Dispatch,
  workingStep: tStep,
  viewingStep: tStep
};

const StepperContainer = ({
                            steps,
                            viewingIndex,
                            jumpTo,
                            workingStep,
                            viewingStep
                          }: StepperLayoutProps) => {
  const classes = makeStyles(styles.stepperContainer)();

  const [viewingStepRef, setViewingStepRef] = useState(null);
  useEffect(() => {
    if (viewingStepRef) {
      // eslint-disable-next-line react/no-find-dom-node
      const node = ReactDOM.findDOMNode(viewingStepRef);
      if (node) {
        node.scrollIntoView({ block: 'center',  behavior: 'smooth' });
      }
    }
  }, [viewingStepRef]);

  return (
    <Stepper
      nonLinear
      activeStep={viewingIndex}
      orientation="vertical"
      className={classes.root}
    >
      {steps.map((s, idx) => {
        const fail = s.status === STEP_STATUS.FAIL;
        const stepProps = {};
        const labelProps = {
          error: fail
        };
        const stepButtonProps = {};
        if (workingStep === s) {
          stepButtonProps.icon = <Loop className={classes.stepIconDoing}/>;
        }
        return (
          <Step key={s.name} {...stepProps} >
            <StepButton
              completed={s.status === STEP_STATUS.FINISHED}
              onClick={() => jumpTo(idx)}
              className={classes.stepButton}
              {...stepButtonProps}
              ref={
                r => {
                  if (s === viewingStep) {
                    setViewingStepRef(r);
                  }
                }
              }
            >
              <StepLabel {...labelProps} >
                {s.name}
              </StepLabel>
            </StepButton>
            <StepContent>
              <Timer step={s}/>
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
