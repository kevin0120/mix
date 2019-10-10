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
import clsx from 'clsx';
import STEP_STATUS from '../../modules/step/constants';
import styles from './styles';
import * as oSel from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { Dispatch } from '../../modules/indexReducer';
import Timer from './Timer';
import type { tClsStep } from '../../modules/step/Step';

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
  steps: Array<tClsStep>,
  viewingIndex: number,
  jumpTo: Dispatch,
  workingStep: tClsStep,
  viewingStep: tClsStep
};

// 步骤条
const StepperContainer = ({
                            steps,
                            viewingIndex,
                            jumpTo,
                            workingStep,
                            viewingStep
                          }: StepperLayoutProps) => {
  const classes = makeStyles(styles.stepperContainer)();

  const viewingRef = useRef(null);
  const viewingNode = viewingRef?.current;
  useEffect(() => {
    if (viewingNode) {
      // eslint-disable-next-line react/no-find-dom-node
      const node: null | Element | Text = ReactDOM.findDOMNode(viewingNode);
      if (node && node.scrollIntoView && typeof node.scrollIntoView === 'function') {
        ((node: any): Element).scrollIntoView({ block: 'center', behavior: 'smooth' });
      }
    }
  }, [viewingNode]);

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
          stepButtonProps.icon = (
            <Loop
              className={clsx(classes.stepIconDoing, {
                [classes.fail]: fail
              })}
            />
          );
        }

        if (s === viewingStep) {
          stepButtonProps.ref = viewingRef;
        } else {
          stepButtonProps.ref = null;
        }

        return (
          <Step key={s.name}>
            <StepButton
              completed={s.status === STEP_STATUS.FINISHED}
              onClick={() => jumpTo(idx)}
              className={classes.stepButton}
              {...stepButtonProps}
              ref={s === viewingStep ? viewingRef : () => {
              }}
            >
              <StepLabel {...labelProps}>
                <Typography variant="h6">{s.name}</Typography>
              </StepLabel>
            </StepButton>
            <StepContent>
              <Timer step={s}/>
              {(Object.keys(s.payload.info || {}) || []).map(k => (
                <div className={classes.infoRow} key={k}>
                  <Typography variant="body1">{k || ''}</Typography>
                  <Typography variant="body1">
                    {s.payload.info && s.payload.info[k] || ''}
                  </Typography>
                </div>
              ))}
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
