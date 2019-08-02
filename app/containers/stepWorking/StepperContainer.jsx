// @flow
import { makeStyles, StepContent, Typography } from '@material-ui/core';
import { connect } from 'react-redux';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import StepLabel from '@material-ui/core/StepLabel';
import React, { useRef, useEffect } from 'react';
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

const viewingRef=useRef(null);

  useEffect(() => {
    if (viewingRef?.current) {
      // eslint-disable-next-line react/no-find-dom-node
      const node = ReactDOM.findDOMNode(viewingRef.current);
      if (node) {
        node.scrollIntoView({ block: 'center', behavior: 'smooth' });
      }
    }
  }, [viewingRef.current]);

  return (
    <Stepper
      nonLinear
      activeStep={viewingIndex}
      orientation="vertical"
      className={classes.root}
    >
      {steps.map((s, idx) => {
        const fail = s.status === STEP_STATUS.FAIL;

        const labelProps = {
          error: fail
        };

        const stepButtonProps = {};

        if (workingStep === s) {
          stepButtonProps.icon = <Loop className={classes.stepIconDoing}/>;
        }

        if(s===viewingStep){
          stepButtonProps.ref=viewingRef;
        }else{
          stepButtonProps.ref=null;
        }

        return (
          <Step key={s.name}>
            <StepButton
              completed={s.status === STEP_STATUS.FINISHED}
              onClick={() => jumpTo(idx)}
              className={classes.stepButton}
              {...stepButtonProps}
              ref={s===viewingStep? viewingRef:()=>{}}
            >
              <StepLabel {...labelProps} >
                <Typography variant="h6">
                {s.name}
                </Typography>
              </StepLabel>
            </StepButton>
            <StepContent>
              <Timer step={s}/>
              <Typography variant="body1">{s.info}</Typography>
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
