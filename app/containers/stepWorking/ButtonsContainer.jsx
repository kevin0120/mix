// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Menu } from '@material-ui/icons';
import Dialog from '@material-ui/core/Dialog';
import Typography from '@material-ui/core/Typography';
import DialogContent from '@material-ui/core/DialogContent';
import Button from '../../components/CustomButtons/Button';
import styles from './styles';
import type { Dispatch } from '../../modules/indexReducer';
import * as orderSelectors from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { tStep, tStepArray } from '../../modules/order/model';

const mapState = (state, props) => ({
  ...props,
  viewingStep: orderSelectors.viewingStep(state.order) || {},
  workingStep: orderSelectors.workingStep(state.order) || {},
  steps: orderSelectors.orderSteps(orderSelectors.viewingOrder(state.order)) || [],
  viewingIndex: orderSelectors.viewingIndex(state.order) || 0,
  isPending: orderSelectors.isPending(orderSelectors.viewingOrder(state.order)),
  isCancel: orderSelectors.isCancel(orderSelectors.viewingOrder(state.order)),
  pendingable: orderSelectors.pendingable(orderSelectors.viewingOrder(state.order)),
  cancelable: orderSelectors.cancelable(orderSelectors.viewingOrder(state.order))
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
  viewingStep: tStep,
  workingStep: tStep,
  viewingIndex: number,
  steps: tStepArray,
  next: Dispatch,
  action: Dispatch,
  previous: Dispatch,
  doNextStep: Dispatch,
  doPreviousStep: Dispatch,
  cancelOrder: Dispatch,
  pendingOrder: Dispatch,
  isPending: boolean,
  isCancel: boolean,
  pendingable: boolean,
  cancelable: boolean
};

const ButtonsContainer = ({
                            viewingStep,
                            workingStep,
                            next,
                            steps,
                            viewingIndex,
                            action,
                            previous,
                            doNextStep,
                            doPreviousStep,
                            cancelOrder,
                            pendingOrder,
                            isPending,
                            isCancel,
                            pendingable,
                            cancelable
                          }: ButtonsContainerProps) => {

  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;


  const [dialogOpen, setDialogOpen] = useState(false);


  return <div className={classes.root}>
    <div>
      {
        isPending || pendingable || cancelable ?
          <React.Fragment>
            <Button
              type="button"
              color="warning"
              onClick={() => setDialogOpen(true)}
            >
              <Menu fontSize="large" className={classes.menuIcon}/>
            </Button>
            <Dialog
              open={dialogOpen}
              onClose={() => setDialogOpen(false)}
            >
              <div style={{ backgroundColor: 'white' }}>
                <DialogContent className={classes.dialogContainer}>
                  {isPending ?
                    <Button
                      type="button"
                      // onClick={() => {
                      //   pendingOrder(true);
                      //   setDialogOpen(false);
                      // }}
                      variant="contained"
                      color="primary"
                      className={classes.bigButton}
                    >
                      <Typography variant="h3">
                        continue
                      </Typography>
                    </Button> :
                    (pendingable && <Button
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
                    </Button>) || null
                  }
                  {
                    cancelable ? <Button
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
                    </Button> : null
                  }
                </DialogContent>
              </div>
            </Dialog>
          </React.Fragment> : null
      }
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
      <Button
        disabled={noNext || viewingStep !== workingStep || !viewingStep?.skippable}
        type="button"
        onClick={() => doNextStep()}
        color="primary"
      >
        {'skip'}
      </Button>
      <Button
        disabled={viewingStep !== workingStep || !viewingStep?.undoable}
        type="button"
        onClick={() => doPreviousStep()}
        color="primary"
      >
        {'undo'}
      </Button>
    </div>
    <div>{action}</div>
  </div>;
};

export default connect(mapState, mapDispatch)(ButtonsContainer);
