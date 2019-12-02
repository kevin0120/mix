// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardMedia from '@material-ui/core/CardMedia';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';
import CardContent from '@material-ui/core/CardContent';
import { sortBy } from 'lodash-es';
import { orderActions } from '../../modules/order/action';
import { ORDER_STATUS } from '../../modules/order/constants';
import {
  todoOrders,
  doneOrders,
  exceptOrders,
  doingOrders
} from '../../modules/order/selector';
import styles from './styles';
import settingImg from '../../../resources/imgs/setting.png';
import type { Dispatch } from '../../modules/typeDef';
import { withI18n } from '../../i18n';
import type { IOrder } from '../../modules/order/interface/IOrder';

type tOP = {||};

type tSP = {|
  ...tOP,
  orderList: Array<IOrder>
|};

type tDP = {|
  view: Dispatch,
  getList: Dispatch
|};

type Props = {
  ...tOP,
  ...tSP,
  ...tDP
};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  orderList: state.order.list
});

const mapDispatch: tDP = {
  view: orderActions.tryView,
  getList: orderActions.getList
};

function HomeOperationList(props: Props) {
  const classes = makeStyles(styles)();
  const { view, orderList, getList } = props;
  const retOrderList = sortBy(
    orderList,
    (o: IOrder) =>
      (o.datePlannedStart && new Date(o.datePlannedStart)) || Date.now()
  );

  const onCardClick = order => {
    view(order);
  };

  const statusMap = {
    [ORDER_STATUS.TODO]: classes.statusTodo,
    [ORDER_STATUS.WIP]: classes.statusWIP,
    [ORDER_STATUS.DONE]: classes.statusDone,
    [ORDER_STATUS.CANCEL]: classes.statusCancel,
    [ORDER_STATUS.PENDING]: classes.statusPending,
    [ORDER_STATUS.FAIL]: classes.statusFail
  };

  useEffect(() => {
    getList();
  }, [getList]);

  const renderOrders = (t, orders, size, title) => (
    <React.Fragment>
      <Grid item xs={12} className={classes.listTitle}>
        <Typography gutterBottom variant="h6" align="left">
          {title}
        </Typography>
      </Grid>
      {/* eslint-disable-next-line react/no-array-index-key */}
      {(orders &&
        orders.map((order: IOrder) =>
          order ? (
            <Grid item xs={size} key={`${order.code}`}>
              <Paper square className={classes.orderCardContainer}>
                <CardActionArea
                  className={classes.orderCard}
                  onClick={() => onCardClick(order)}
                >
                  <div
                    className={clsx(
                      statusMap[order.status || ORDER_STATUS.TODO],
                      classes.statusIndicator
                    )}
                  />
                  <CardMedia
                    className={classes.image}
                    src={order.productTypeImage || settingImg}
                    component="img"
                    style={{ overflow: 'hidden' }}
                    alt={order.productCode}
                  />
                  <CardContent className={classes.info}>
                    <Typography
                      variant="body1"
                      align="left"
                      className={classes.orderNameText}
                    >
                      {order.code}
                    </Typography>
                    {[
                      order.productCode,
                      order.workcenter,
                      order.datePlannedStart &&
                      order.datePlannedStart.toLocaleString()
                    ].map((d, idx) => (
                      <Typography
                        key={`${d}-${idx}`}
                        variant="body2"
                        color="textSecondary"
                        align="left"
                        className={classes.orderInfoText}
                      >
                        {d}
                      </Typography>
                    ))}
                    <Typography
                      variant="body2"
                      color="textSecondary"
                      align="left"
                      className={classes.orderStatusText}
                    >
                      {t(`OrderStatus.${order.status || ORDER_STATUS.TODO}`)}
                    </Typography>
                  </CardContent>
                </CardActionArea>
              </Paper>
            </Grid>
          ) : null
        )) ||
      null}
    </React.Fragment>
  );

  return withI18n(
    t => (
      <div className={classes.root}>
        <Grid
          container
          className={clsx(classes.container, classes.bgEven)}
          justify="center"
          spacing={4}
        >
          <Grid
            item
            container
            xs={6}
            spacing={1}
            alignItems="flex-start"
            alignContent="flex-start"
            justify="flex-start"
            direction="row"
            className={classes.bgOdd}
          >
            {renderOrders(
              t,
              [...doingOrders(retOrderList), ...todoOrders(retOrderList)],
              6,
              t(`OrderStatus.wip`)
            )}
          </Grid>
          <Grid
            item
            container
            xs={3}
            spacing={1}
            alignItems="flex-start"
            alignContent="flex-start"
            justify="flex-start"
            direction="row"
            className={classes.bgEven}
          >
            {renderOrders(
              t,
              doneOrders(retOrderList),
              12,
              t('OrderStatus.done')
            )}
          </Grid>
          <Grid
            item
            container
            xs={3}
            spacing={1}
            alignItems="flex-start"
            alignContent="flex-start"
            justify="flex-start"
            direction="row"
            className={classes.bgOdd}
          >
            {renderOrders(
              t,
              exceptOrders(retOrderList),
              12,
              t('OrderStatus.excp')
            )}
          </Grid>
        </Grid>
      </div>
    ),
    'translations'
  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(HomeOperationList);
