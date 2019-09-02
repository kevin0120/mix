// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper } from '@material-ui/core';
import clsx from 'clsx';
import { I18n } from 'react-i18next';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';
import type { tOrderStatus } from '../../modules/order/model';
import { ORDER_STATUS } from '../../modules/order/model';
import logo from '../../../resources/imgs/logo.jpg';
import { translation as trans, stepWorkingNS } from './local';
import { withI18n } from '../../i18n';

type Props = {
  status: ?tOrderStatus,
  name: ?string
};

const statusMap = classes => ({
  empty: null,
  [ORDER_STATUS.TODO]: classes.statusTodo,
  [ORDER_STATUS.WIP]: classes.statusWIP,
  [ORDER_STATUS.DONE]: classes.statusDone,
  [ORDER_STATUS.CANCEL]: classes.statusCancel,
  [ORDER_STATUS.PENDING]: classes.statusPending
});

function StepWorking({ status, name, desc }: Props) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);
  const [description, bindDescription] = useState(null);
  return withI18n(t => (
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
            [{status ? t(`OrderStatus.${status}`, { ns: 'translations' }) : t('notSelected')}]
          </Typography>
          {/* <Typography variant="h5">{name || ''}</Typography> */}
          {/* <Typography variant="h5">{desc || ''}</Typography> */}
          {
            desc && desc.split('\t\t').map(d =>
              <React.Fragment key={d}>
                <Typography variant="h5">
                  {d || ''}
                </Typography>
                <span style={{ width: '40px' }}/>
              </React.Fragment>
            )
          }
        </div>
        <img alt="logo" src={logo} className={classes.logo}/>
      </Paper>
      <div className={classes.main}>
        <Paper square classes={{ root: classes.leftContainer }}>
          <ButtonsContainer action={action}/>
          <StepPageContainer
            bindAction={bindAction}
            bindDescription={bindDescription}
            description={description}/>
        </Paper>
        <div className={classes.rightContainer}>
          <Paper square className={classes.stepperContainer}>
            <StepperContainer/>
          </Paper>
        </div>
      </div>
    </div>
  ), stepWorkingNS);
}

const mapState = (state, props) => {
  const vOrder = orderSelectors.viewingOrder(state.order);
  return {
    ...props,
    status: vOrder?.status,
    name: vOrder?.name,
    desc: vOrder?.desc
  };
};

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
