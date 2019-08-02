// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper } from '@material-ui/core';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';
import type { tOrder } from '../../modules/order/model';
import { I18n } from 'react-i18next';
import { ORDER_STATUS } from '../../modules/order/model';
import clsx from 'clsx';

type Props = {
  viewingOrder: tOrder
};

const statusMap = (classes) => ({
  [ORDER_STATUS.TODO]: classes.statusTodo,
  [ORDER_STATUS.WIP]: classes.statusWIP,
  [ORDER_STATUS.DONE]: classes.statusDone,
  [ORDER_STATUS.CANCEL]: classes.statusCancel,
  [ORDER_STATUS.PENDING]: classes.statusPending
});

function StepWorking({ viewingOrder }: Props) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  return (
    <I18n ns="translations">
      {t => (
        <div className={classes.root}>
          <Paper square className={classes.orderInfoContainer}>
            <Typography variant="h5" className={clsx(statusMap(classes)[viewingOrder?.status],classes.orderStatus)}>
              [{viewingOrder?t(viewingOrder.status) || '':'未选中工单'}]
            </Typography>
            <Typography variant="h5">{viewingOrder?.name || ''}</Typography>
          </Paper>
          <div className={classes.main}>
            <Paper
              square
              classes={{ root: classes.leftContainer }}
            >
              <ButtonsContainer action={action}/>
              <StepPageContainer bindAction={bindAction}/>
            </Paper>
            <div className={classes.rightContainer}>
              <Paper square className={classes.stepperContainer}>
                <StepperContainer/>
              </Paper>
            </div>


          </div>
        </div>
      )}

    </I18n>
  );
}

const mapState = (state, props) => ({
  ...props,
  viewingOrder: orderSelectors.viewingOrder(state.order) || null
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
