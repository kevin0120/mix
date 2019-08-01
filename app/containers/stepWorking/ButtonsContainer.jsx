// @flow

import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Menu } from '@material-ui/icons';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import { I18n } from 'react-i18next';
import Button from '../../components/CustomButtons/Button';
import styles from './styles';
import actionStepWorkingDef from './type';
import type { Dispatch } from '../../modules/indexReducer';
import * as oSel from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { tOrder, tStep, tStepArray } from '../../modules/order/model';

const mapState = (state, props) => {
  const vOrder = oSel.viewingOrder(state.order);
  return {
    ...props,
    viewingOrder: vOrder,
    viewingStep: oSel.viewingStep(state.order) || {},
    workingStep: oSel.workingStep(oSel.workingOrder(state.order)) || {},
    steps: oSel.orderSteps(vOrder) || [],
    viewingIndex: oSel.viewingIndex(state.order) || 0,
    isPending: oSel.isPending(vOrder),
    isCancel: oSel.isCancel(vOrder),
    pendingable: oSel.pendingable(vOrder),
    cancelable: oSel.cancelable(vOrder)
  };
};

const mapDispatch = {
  next: orderActions.nextStep,
  doNextStep: orderActions.doNextStep,
  previous: orderActions.previousStep,
  doPreviousStep: orderActions.doPreviousStep,
  cancelOrder: orderActions.cancelOrder,
  pendingOrder: orderActions.pendingOrder,
  workOn: orderActions.workOn
};

type ButtonsContainerProps = {
  viewingOrder: tOrder,
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
  // isCancel: boolean,
  pendingable: boolean,
  cancelable: boolean,
  workOn: Dispatch
};

const ButtonsContainer = ({
                            viewingOrder,
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
                            pendingable,
                            cancelable,
                            workOn
                          }: ButtonsContainerProps) => {

  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;


  const [dialogOpen, setDialogOpen] = useState(false);


  return <I18n ns="translations">
    {t => (<div className={classes.root}>
    <div>
      {
        isPending || pendingable || cancelable ?
          <React.Fragment>
            <Button
              type="button"
              color="warning"
              onClick={() => setDialogOpen(true)}
            >
                <Menu fontSize="inherit" className={classes.menuIcon}/>
              {' '}
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
                      onClick={() => {
                        workOn(viewingOrder);
                        setDialogOpen(false);
                      }}
                      variant="contained"
                      color="primary"
                      className={classes.bigButton}
                    >
                        continue
                    </Button> :
                    (pendingable && <Button
                      type="button"
                      onClick={() => {
                        pendingOrder(viewingOrder);
                        setDialogOpen(false);
                      }}
                      variant="contained"
                      color="warning"
                      className={classes.bigButton}
                    >
                        pending
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
                        cancel
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
        disabled={viewingStep !== workingStep || !viewingStep?.skippable}
        type="button"
        onClick={() => doNextStep()}
        color="primary"
      >
        {t(actionStepWorkingDef.SKIP)}
      </Button>
      <Button
        disabled={viewingStep !== workingStep || !viewingStep?.undoable}
        type="button"
        onClick={() => doPreviousStep()}
        color="primary"
      >
        {t(actionStepWorkingDef.UNDO)}
      </Button>
    </div>
    <div>{action}</div>
  </div>)}
  </I18n>
};

export default connect(mapState, mapDispatch)(ButtonsContainer);
