// @flow

import React, { useState } from 'react';
import type { Node} from 'react'
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Menu } from '@material-ui/icons';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import { I18n } from 'react-i18next';
import Button from '../../components/CustomButtons/Button';
import styles from './styles';
import { translation as trans, stepWorkingNS } from './local';
import * as oSel from '../../modules/order/selector';
import { orderActions } from '../../modules/order/action';
import type { tStep, tStepArray } from '../../modules/order/model';
import dialogActions from '../../modules/dialog/action';
import { tNS, withI18n } from '../../i18n';
import Table from '../../components/Table/Table';
import modelViewerActions from '../../modules/modelViewer/action';
import type { tClsOrder } from '../../modules/order/Order';
import type { orderTriggerType, updateStateActionType } from '../../modules/order/action';

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
  workOn: orderActions.workOn,
  viewModelDialog: dialogActions.dialogShow,
  viewModel: modelViewerActions.open
};

type ButtonsContainerProps = {
  viewingOrder: tClsOrder,
  viewingIndex: number,
  viewingStep: tStep,
  workingStep: tStep,
  viewingIndex: number,
  steps: tStepArray,
  /* eslint-disable flowtype/no-weak-types */
  next: ()=> any,
  action: Node,
  previous: ()=> any,
  doNextStep: ()=> any,
  doPreviousStep: ()=> any,
  cancelOrder: (order: tClsOrder)=> updateStateActionType,
  pendingOrder: (order: tClsOrder)=> updateStateActionType,
  isPending: boolean,
  // isCancel: boolean,
  pendingable: boolean,
  cancelable: boolean,
  workOn: (order: tClsOrder)=> orderTriggerType,
  viewModel: any, // 查看的三维模型
  viewModelDialog: any
  /* eslint-enable flowtype/no-weak-types */
};

const ButtonsContainer: (ButtonsContainerProps) => Node = ({
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
                            workOn,
                            viewModel,
                            viewModelDialog
                          }: ButtonsContainerProps) => {

  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;


  const [dialogOpen, setDialogOpen] = useState(false);

  const modelsData = (viewingOrder?.payload?.models &&
    viewingOrder.payload.models.map((m) => [
      m.name,
      m.desc,
      withI18n(t => (
        <Button
          color="primary"
          type="button"
          onClick={() => viewModel(m.url)}
        >{t(trans.view)}</Button>
      ), stepWorkingNS)
    ])) || [];
  const modelsTableDialog = {
    buttons: [
      {
        label: 'Common.Close',
        color: 'warning'
      }
    ],
    title: tNS(trans.viewModel, stepWorkingNS),
    content: (
      withI18n(t => (
        <Table
          tableHeaderColor="info"
          tableHead={[
            t(trans.name),
            t(trans.desc),
            t(trans.action)
          ]}
          tableData={modelsData}
          colorsColls={['info']}
        />
      ), stepWorkingNS)
    )
  };

  return withI18n(t => (
    <div className={classes.root}>
      <div>
        {
          isPending || pendingable || cancelable ?
            <React.Fragment>
              <Button
                justIcon
                type="button"
                color="github"
                onClick={() => setDialogOpen(true)}
              >
                <Menu fontSize="inherit" className={classes.menuIcon}/>
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
                        {t(trans.continueDoing)}
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
                        {t(trans.pending)}
                      </Button>) || null
                    }
                    {
                      cancelable ? <Button
                        type="button"
                        color="danger"
                        className={classes.bigButton}
                        onClick={() => {
                          cancelOrder(viewingOrder);
                          setDialogOpen(false);
                        }}
                      >
                        {t(trans.cancel)}
                      </Button> : null
                    }
                  </DialogContent>
                </div>
              </Dialog>
            </React.Fragment> : null
        }
        <Button
          color="primary"
          type="button"
          onClick={() => viewModelDialog(modelsTableDialog)}
        >
          {t(trans.viewModel)}
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
        <Button
          disabled={viewingStep !== workingStep || !viewingStep?.skippable}
          type="button"
          onClick={() => doNextStep()}
          color="tumblr"
        >
          {t(trans.skip)}
        </Button>
        <Button
          disabled={viewingStep !== workingStep || !viewingStep?.undoable}
          type="button"
          onClick={() => doPreviousStep()}
          color="danger"
        >
          {t(trans.undo)}
        </Button>
      </div>
      <div>{action}</div>
    </div>), stepWorkingNS);
};

export default connect<ButtonsContainerProps,*, _,_,_,_>(mapState, mapDispatch)(ButtonsContainer);
