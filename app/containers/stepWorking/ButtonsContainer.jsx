// @flow

import type { Node } from 'react';
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Menu } from '@material-ui/icons';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import Button from '../../components/CustomButtons/Button';
import styles from './styles';
import { stepWorkingNS, translation as trans } from './local';
import * as oSel from '../../modules/order/selector';
import type { tActOrderTrigger, tActUpdateState } from '../../modules/order/action';
import { orderActions } from '../../modules/order/action';
import dialogActions from '../../modules/dialog/action';
import { tNS, withI18n } from '../../i18n';
import Table from '../../components/Table/Table';
import modelViewerActions from '../../modules/modelViewer/action';
import type { IOrder } from '../../modules/order/interface/IOrder';
import type { IWorkStep } from '../../modules/step/interface/IWorkStep';
import PDFViewer from '../../components/PDFViewer';
import { CommonLog, defaultClient } from '../../common/utils';
import { BlockReasonDialog } from '../../components/BlockReasonDialog';
import STEP_STATUS from '../../modules/step/constants';
import { Typography } from '@material-ui/core';
import { ORDER_STATUS } from '../../modules/order/constants';

const mapState = (state, props) => {
  const vOrder = oSel.viewingOrder(state.order);
  const wOrder = oSel.workingOrder(state.order);
  return {
    ...props,
    viewingOrder: vOrder,
    workingOrder: wOrder,
    canDoAnotherStep: state.setting.systemSettings.canDoAnotherStep || false,
    viewingStep: oSel.viewingStep(state.order) || {},
    workingStep: oSel.workingStep(oSel.workingOrder(state.order)) || {},
    steps: oSel.orderSteps(vOrder) || [],
    viewingIndex: oSel.viewingIndex(state.order) || 0,
    isPending: oSel.isPending(vOrder),
    isCancel: oSel.isCancel(vOrder),
    pendingable: oSel.pendingable(vOrder),
    cancelable: oSel.cancelable(vOrder),
    canReportFinish: oSel.canReportFinish(vOrder) || false,
    reportFinishEnabled: state.setting.systemSettings.reportFinish,
    blockReasons: state.order.blockReasons || [],
    canRedoOrders: state.setting?.systemSettings?.canRedoOrders &&
      !wOrder && (vOrder?.status === ORDER_STATUS.DONE || vOrder?.status === ORDER_STATUS.FAIL || vOrder?.status===ORDER_STATUS.CANCEL)
  };
};

const mapDispatch = {
  next: orderActions.nextStep,
  finishStep: orderActions.finishStep,
  previous: orderActions.previousStep,
  doPreviousStep: orderActions.doPreviousStep,
  doAnotherStep: orderActions.doAnotherStep,
  cancelOrder: orderActions.cancelOrder,
  pendingOrder: orderActions.pendingOrder,
  tryWorkOn: orderActions.tryWorkOn,
  showDialog: dialogActions.dialogShow,
  viewModel: modelViewerActions.open,
  reportFinish: orderActions.reportFinish,
  redoOrder: orderActions.redoOrder
};

/* eslint-disable flowtype/no-weak-types */
type ButtonsContainerProps = {
  viewingOrder: IOrder,
  workingOrder: IOrder,
  viewingIndex: number,
  viewingStep: IWorkStep,
  workingStep: IWorkStep,
  steps: Array<IWorkStep>,
  next: () => any,
  action: Node,
  previous: () => any,
  finishStep: IWorkStep => any,
  doPreviousStep: () => any,
  doAnotherStep: () => any,
  cancelOrder: (order: IOrder) => tActUpdateState,
  pendingOrder: (order: IOrder) => tActUpdateState,
  isPending: boolean,
  // isCancel: boolean,
  pendingable: boolean,
  cancelable: boolean,
  tryWorkOn: (order: IOrder) => tActOrderTrigger,
  redoOrder: (order: IOrder) => tActOrderTrigger,
  viewModel: any, // 查看的三维模型
  viewModelDialog: any,
  canRedoOrders: boolean
};
/* eslint-enable flowtype/no-weak-types */

const ButtonsContainer: ButtonsContainerProps => Node = ({
  viewingOrder,
  workingOrder,
  canDoAnotherStep,
  viewingStep,
  workingStep,
  next,
  steps,
  viewingIndex,
  action,
  previous,
  finishStep,
  doPreviousStep,
  doAnotherStep,
  cancelOrder,
  pendingOrder,
  isPending,
  pendingable,
  cancelable,
  tryWorkOn,
  viewModel,
  showDialog,
  reportFinish,
  canReportFinish,
  reportFinishEnabled,
  blockReasons,
  canRedoOrders, // 工单是否能重新作业
  redoOrder
}: ButtonsContainerProps) => {
  const classes = makeStyles(styles.buttonsContainer)();
  const noPrevious = steps.length <= 0 || viewingIndex <= 0;
  const noNext = steps.length <= 0 || viewingIndex >= steps.length - 1;

  const [dialogOpen, setDialogOpen] = useState(false);

  const modelsData =
    (viewingOrder &&
      viewingOrder.payload &&
      viewingOrder.payload.products &&
      viewingOrder.payload.products.map(m => [
        m.code,
        m.desc,
        withI18n(
          t => (
            <Button
              color="primary"
              type="button"
              onClick={() => viewModel(m.url)}
            >
              {t(trans.view)}
            </Button>
          ),
          stepWorkingNS
        )
      ])) ||
    [];

  const [pdfUrl, setPdfUrl] = useState('');

  const url = viewingOrder?.payload?.worksheet?.url;
  useEffect(() => {
    if (!url) {
      return;
    }
    defaultClient
      .get(url)
      .then(resp => {
        setPdfUrl(resp?.request?._redirectable?._currentUrl || '');
      }).catch((err) => {
      CommonLog.lError(err);
    });
  }, [url]);

  const modelsTableDialog = {
    maxWidth: 'md',
    buttons: [
      {
        label: 'Common.Close',
        color: 'warning'
      }
    ],
    title: tNS(trans.viewModel, stepWorkingNS),
    content: withI18n(
      t => (
        <Table
          tableHeaderColor="info"
          tableHead={[t(trans.name), t(trans.desc), t(trans.action)]}
          tableData={modelsData}
          colorsColls={['info']}
        />
      ),
      stepWorkingNS
    )
  };
  const fileDialog = {
    maxWidth: 'lg',
    buttons: [{
      label: 'Common.Close',
      color: 'warning'
    }],
    title: tNS(trans.viewFile, stepWorkingNS),
    content: <PDFViewer file={pdfUrl || url}/>
  };

  const envDialog = {
    maxWidth: 'xl',
    buttons: [{
      label: 'Common.Close',
      color: 'warning'
    }],
    content: <iframe src="http://172.26.214.80:8091/a/login" style={{ width: '80vw', height: '75vh' }}/>
  };

  return withI18n(
    t => (
      <div className={classes.root}>
        <div>
          {isPending || pendingable || cancelable ? (
            <React.Fragment>
              <Button
                justIcon
                type="button"
                color="github"
                onClick={() => setDialogOpen(true)}
              >
                <Menu fontSize="inherit" className={classes.menuIcon}/>
              </Button>
              <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
                <div style={{ backgroundColor: 'white' }}>
                  <DialogContent className={classes.dialogContainer}>
                    {isPending ? (
                      <Button
                        type="button"
                        onClick={() => {
                          tryWorkOn(viewingOrder);
                          setDialogOpen(false);
                        }}
                        variant="contained"
                        color="primary"
                        className={classes.bigButton}
                      >
                        <Typography
                          variant="h4"
                          color="inherit"
                          align="left"
                          className={classes.orderInfoText}
                        >
                          {t(trans.continueDoing)}
                        </Typography>
                      </Button>
                    ) : (
                      (pendingable && (
                        <BlockReasonDialog
                          blockReasons={blockReasons}
                          AnchorButton={({ onClick }) => <Button
                            type="button"
                            onClick={onClick}
                            variant="contained"
                            color="warning"
                            className={classes.bigButton}
                          >
                            <Typography
                              variant="h4"
                              color="inherit"
                              align="left"
                              className={classes.orderInfoText}
                            >
                              {t(trans.pending)}
                            </Typography>
                          </Button>}
                          onConfirm={(reason) => {
                            pendingOrder(viewingOrder, reason);
                            setDialogOpen(false);
                          }}
                        />
                      )) || null
                    )}
                    {cancelable ? (
                      <Button
                        type="button"
                        color="danger"
                        className={classes.bigButton}
                        onClick={() => {
                          cancelOrder(viewingOrder);
                          setDialogOpen(false);
                        }}
                      >
                        <Typography
                          variant="h4"
                          color="inherit"
                          align="left"
                          className={classes.orderInfoText}
                        >
                          {t(trans.cancel)}
                        </Typography>
                      </Button>
                    ) : null}
                  </DialogContent>
                </div>
              </Dialog>
            </React.Fragment>
          ) : null}
          {canReportFinish && reportFinishEnabled ? <Button
            color="warning"
            type="button"
            onClick={() => reportFinish(viewingOrder)}
          >
            {t(trans.reportFinish)}
          </Button> : null}
          <Button
            color="primary"
            type="button"
            onClick={() => showDialog(modelsTableDialog)}
          >
            {t(trans.viewModel)}
          </Button>
          <Button
            color="primary"
            type="button"
            onClick={() => showDialog(fileDialog)}
          >
            {t(trans.viewFile)}
          </Button>
          <Button
            color="primary"
            type="button"
            disabled={!viewingOrder}
            onClick={() => showDialog(envDialog)}
          >
            {'生产指导'}
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
            onClick={() => finishStep(viewingStep)}
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
          {
            canRedoOrders ? <Button
              disabled={!canRedoOrders}
              type="button"
              onClick={() => redoOrder(viewingOrder)}
              color="danger"
            >
              {t(trans.redoOrder)}
            </Button> : null
          }
          {canDoAnotherStep && viewingOrder === workingOrder ? <Button
            type="button"
            color="info"
            disabled={viewingStep === workingStep || viewingStep.status === STEP_STATUS.FINISHED}
            onClick={() => {
              doAnotherStep(viewingStep);
            }}>
            {t(trans.startViewing)}
          </Button> : null}

        </div>
        <div>{action}</div>
      </div>
    ),
    stepWorkingNS
  );
};

export default connect<ButtonsContainerProps, *, _, _, _, _>(
  mapState,
  mapDispatch
)(ButtonsContainer);
