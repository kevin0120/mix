import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/es/Button/Button';
import styles from './styles';
import type { Dispatch } from '../../modules/indexReducer';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';

const mapState = (state, props) => ({
  ...props,
  viewingStep: orderSelectors.viewingStep(state.order) || {},
  processingStep: orderSelectors.processingStep(state.order) || {},
  steps: orderSelectors.orderSteps(state.order) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0,
});

const mapDispatch = {
  next: orderActions.nextStep,
  doNextStep: orderActions.doNextStep,
  previous: orderActions.previousStep,
  doPreviousStep: orderActions.doPreviousStep
};

type ButtonsContainerProps = {
  viewingIndex: number,
  viewingStep: {},
  processingStep: {},
  viewingIndex: number,
  steps: [],
  next: Dispatch,
  action: Dispatch,
  previous: Dispatch,
  doNextStep: Dispatch,
  doPreviousStep: Dispatch
};

const ButtonsContainer = ({
                            viewingStep,
                            processingStep,
                            next,
                            steps,
                            viewingIndex,
                            action,
                            previous,
                            doNextStep,
                            doPreviousStep,
                          }: ButtonsContainerProps) => {
  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;

  return <div className={classes.root}>
    <div>
      <Button
        disabled={noPrevious}
        type="button"
        onClick={() => previous()}
      >
        {'<<'}
      </Button>
      <Button disabled={noNext} type="button" onClick={() => next()}>
        {'>>'}
      </Button>
      {viewingStep?.skippable && (
        <Button
          disabled={noNext || viewingStep !== processingStep}
          type="button"
          onClick={() => doNextStep()}
        >
          {'skip'}
        </Button>
      )}
      {viewingStep?.undoable && (
        <Button
          disabled={viewingStep !== processingStep}
          type="button"
          onClick={() => doPreviousStep()}
        >
          {'undo'}
        </Button>
      )}
    </div>
    <div>{action}</div>
  </div>;
};

export default connect(mapState,mapDispatch)(ButtonsContainer);
