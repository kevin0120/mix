// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import Button from '../../components/CustomButtons/Button';
import Dialog from '@material-ui/core/Dialog';
import Typography from '@material-ui/core/Typography';
import DialogContent from '@material-ui/core/DialogContent';
import styles from './styles';
import type { Dispatch } from '../../modules/indexReducer';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import {Menu} from '@material-ui/icons';

const mapState = (state, props) => ({
  ...props,
  viewingStep: orderSelectors.viewingStep(state.order) || {},
  processingStep: orderSelectors.processingStep(state.order) || {},
  steps: orderSelectors.orderSteps(state.order) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0
});

const mapDispatch = {
  next: orderActions.nextStep,
  doNextStep: orderActions.doNextStep,
  previous: orderActions.previousStep,
  doPreviousStep: orderActions.doPreviousStep,
  cancelOrder: orderActions.cancelOrder,
  pendingOrder: orderActions.pendingOrder
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
  doPreviousStep: Dispatch,
  cancelOrder: Dispatch,
  pendingOrder: Dispatch
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
                            cancelOrder,
                            pendingOrder
                          }: ButtonsContainerProps) => {

  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;


  const [dialogOpen, setDialogOpen] = useState(false);


  return <div className={classes.root}>
    <Dialog
      open={dialogOpen}
      onClose={() => setDialogOpen(false)}
    >
      <div style={{ backgroundColor: 'white' }}>
        <DialogContent className={classes.dialogContainer}>
          <Button
            type="button"
            onClick={() => {
              pendingOrder(true);
              setDialogOpen(false);
            }}
            variant="contained"
            color="warning"
            className={classes.bigButton}
          >
            <Typography variant="h3">
              pending
            </Typography>
          </Button>
          <Button
            type="button"
            color="danger"
            className={classes.bigButton}
            onClick={() => {
              cancelOrder(true);
              setDialogOpen(false);
            }}
          >
            <Typography variant="h3">
              cancel
            </Typography>
          </Button>
        </DialogContent>
      </div>
    </Dialog>
    <div>
      <Button
        type="button"
        color="warning"
        onClick={() => setDialogOpen(true)}
      >
        <Menu fontSize="large" className={classes.menuIcon}/>
      </Button>
      <Button
        color="primary"
        disabled={noPrevious}
        type="button"
        onClick={() => previous()}
      >
        {'<<'}
      </Button>
      <Button
        disabled={noNext}
        type="button"
        onClick={() => next()}
        color="primary"
      >
        {'>>'}
      </Button>
      {viewingStep?.skippable && (
        <Button
          disabled={noNext || viewingStep !== processingStep}
          type="button"
          onClick={() => doNextStep()}
          color="primary"
        >
          {'skip'}
        </Button>
      )}
      {viewingStep?.undoable && (
        <Button
          disabled={viewingStep !== processingStep}
          type="button"
          onClick={() => doPreviousStep()}
          color="primary"
        >
          {'undo'}
        </Button>
      )}
    </div>
    <div>{action}</div>
  </div>;
};

export default connect(mapState, mapDispatch)(ButtonsContainer);
