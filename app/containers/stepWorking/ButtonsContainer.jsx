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
    cancelable: oSel.cancelable(vOrder),
    canReportFinish: oSel.canReportFinish(vOrder) || false,
    reportFinishEnabled: state.setting.systemSettings.reportFinish,
    blockReasons: state.order.blockReasons || []
  };
};

const mapDispatch = {
  next: orderActions.nextStep,
  finishStep: orderActions.finishStep,
  previous: orderActions.previousStep,
  doPreviousStep: orderActions.doPreviousStep,
  cancelOrder: orderActions.cancelOrder,
  pendingOrder: orderActions.pendingOrder,
  tryWorkOn: orderActions.tryWorkOn,
  showDialog: dialogActions.dialogShow,
  viewModel: modelViewerActions.open,
  reportFinish: orderActions.reportFinish
};

/* eslint-disable flowtype/no-weak-types */
type ButtonsContainerProps = {
  viewingOrder: IOrder,
  viewingIndex: number,
  viewingStep: IWorkStep,
  workingStep: IWorkStep,
  steps: Array<IWorkStep>,
  next: () => any,
  action: Node,
  previous: () => any,
  finishStep: IWorkStep => any,
  doPreviousStep: () => any,
  cancelOrder: (order: IOrder) => tActUpdateState,
  pendingOrder: (order: IOrder) => tActUpdateState,
  isPending: boolean,
  // isCancel: boolean,
  pendingable: boolean,
  cancelable: boolean,
  tryWorkOn: (order: IOrder) => tActOrderTrigger,
  viewModel: any, // 查看的三维模型
  viewModelDialog: any
};
/* eslint-enable flowtype/no-weak-types */

const ButtonsContainer: ButtonsContainerProps => Node = ({
                                                           viewingOrder,
                                                           viewingStep,
                                                           workingStep,
                                                           next,
                                                           steps,
                                                           viewingIndex,
                                                           action,
                                                           previous,
                                                           finishStep,
                                                           doPreviousStep,
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
                                                           blockReasons
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
    defaultClient
      .get(url)
      .then(resp => {
        setPdfUrl(resp?.request?._redirectable?._currentUrl || '');
      }).catch((err)=>{
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
                        {t(trans.continueDoing)}
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
                            {t(trans.pending)}
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
                        {t(trans.cancel)}
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
