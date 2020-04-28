// @flow
import type { Node } from 'react';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Paper, Typography } from '@material-ui/core';
import clsx from 'clsx';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';
import type { tOrderStatus } from '../../modules/order/interface/typeDef';
import { ORDER_STATUS } from '../../modules/order/constants';
import logo from '../../../resources/imgs/logo.jpg';
import { stepWorkingNS } from './local';
import { withI18n } from '../../i18n';
import OrderInfo from './OrderInfo';

const { OrderInfoLeft, OrderInfoRight } = OrderInfo;
const mapState = (state, props) => {
  const vOrder = orderSelectors.viewingOrder(state.order);
  return {
    ...props,
    status: vOrder?.status,
    code: vOrder?.code
  };
};

const mapDispatch = {};

type Props = {
  status: ?tOrderStatus,
  name: ?string,
  code: ?string,
  clickPoint: () => void
};

const statusMap = classes => ({
  empty: null,
  [ORDER_STATUS.TODO]: classes.statusTodo,
  [ORDER_STATUS.WIP]: classes.statusWIP,
  [ORDER_STATUS.DONE]: classes.statusDone,
  [ORDER_STATUS.CANCEL]: classes.statusCancel,
  [ORDER_STATUS.PENDING]: classes.statusPending,
  [ORDER_STATUS.FAIL]: classes.statusFail
});

function StepWorking({ status, code, clickPoint }: Props): Node {
  const classes = makeStyles(styles.layout)();

  type tNodeHook = [Node, ((Node => Node) | Node) => void];
  const [action, bindAction]: tNodeHook = useState(null);
  const [description, bindDescription]: tNodeHook = useState(null);

  return withI18n(
    t => (
      <div className={classes.root}>
        <Paper square className={classes.topBarContainer}>
          <div className={classes.orderInfoContainer}>
            <Typography
              variant="h5"
              className={clsx(
                statusMap(classes)[status || 'empty'],
                classes.orderStatus
              )}
            >
              [
              {status
                ? t(`OrderStatus.${status}`, { ns: 'translations' })
                : t('notSelected')}
              ]
            </Typography>
            <Typography variant="h5">{code || ''}</Typography>
          </div>
          {logo ? <img alt="" src={logo} className={classes.logo}/> : null}
        </Paper>
        <div className={classes.main}>
          <div className={classes.leftContainer}>
            <OrderInfoLeft/>
          </div>
          <Paper square classes={{ root: classes.centerContainer }}>
            <ButtonsContainer action={action}/>
            <StepPageContainer
              bindParentAction={bindAction}
              bindParentDescription={bindDescription}
              description={description}
            />
          </Paper>
          <div className={classes.rightContainer}>
            <Paper square className={classes.stepperContainer}>
              <StepperContainer/>
            </Paper>
          </div>
          <div className={classes.leftContainer}>
            <OrderInfoRight/>
          </div>
        </div>
      </div>
    ),
    stepWorkingNS
  );
}

export default connect<Props, *, _, _, _, _>(
  mapState,
  mapDispatch
)(StepWorking);
