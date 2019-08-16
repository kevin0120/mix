// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardMedia from '@material-ui/core/CardMedia';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';
import CardContent from '@material-ui/core/CardContent';
import { I18n } from 'react-i18next';
import { sortBy } from 'lodash-es';
import { orderActions } from '../../modules/order/action';
import { ORDER_STATUS } from '../../modules/order/model';
import { todoOrders, doneOrders, exceptOrders, doingOrders } from '../../modules/order/selector';
import styles from './styles';
import settingImg from '../../../resources/imgs/setting.png';
import type { Dispatch } from '../../modules/indexReducer';
import type { tOrder } from '../../modules/order/model';


type Props = {
  classes: {},
  orderList: Array<tOrder>,
  getDetail: Dispatch,
  view: Dispatch,
  doPush: Dispatch,
  getList: Dispatch
};

function HomeOperationList(props: Props) {
  const classes = makeStyles(styles)();
  const { view, getDetail, doPush, orderList, getList } = props;

  const retOrderList = sortBy(orderList, (o: tOrder) => new Date(o.plannedDateTime) || Date.now());

  const onCardClick = (order) => {
    getDetail(order);
    view(order);
    doPush('/app/working');
  };

  const statusMap = {
    [ORDER_STATUS.TODO]: classes.statusTodo,
    [ORDER_STATUS.WIP]: classes.statusWIP,
    [ORDER_STATUS.DONE]: classes.statusDone,
    [ORDER_STATUS.CANCEL]: classes.statusCancel,
    [ORDER_STATUS.PENDING]: classes.statusPending
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
      {orders && orders.map((order: tOrder, idx: number) => order ? (<Grid item xs={size} key={`${order.name}${idx}`}>
          <Paper square className={classes.orderCardContainer}>
            <CardActionArea className={classes.orderCard} onClick={() => onCardClick(order)}>
              <div className={clsx(statusMap[order.status || ORDER_STATUS.TODO], classes.statusIndicator)}/>
              <CardMedia
                className={classes.image}
                src={settingImg}
                component="img"
              />
              <CardContent className={classes.info}>
                <Typography variant="body1" align="left" className={classes.orderNameText}>
                  {order.name}
                </Typography>
                <Typography variant="body2" color="textSecondary" align="left" className={classes.orderInfoText}>
                  {order.info}
                </Typography>
                <Typography variant="body2" color="textSecondary" align="left" className={classes.orderStatusText}>
                  {t(order.status || ORDER_STATUS.TODO)}
                </Typography>
              </CardContent>
            </CardActionArea>
          </Paper>
        </Grid>) : null
      ) || null}
    </React.Fragment>
  );

  return <I18n ns="translations">
    {t => (
      <div className={classes.root}>
        <Grid container className={clsx(classes.container, classes.bgEven)} justify="center" spacing={4}>
          <Grid item container xs={6} spacing={1} alignItems="flex-start" alignContent="flex-start"
                justify="flex-start" direction="row" className={classes.bgOdd}>
            {renderOrders(t, [...doingOrders(retOrderList), ...todoOrders(retOrderList)], 6, t('OrderStatus.WIP'))}
          </Grid>
          <Grid item container xs={3} spacing={1} alignItems="flex-start" alignContent="flex-start"
                justify="flex-start" direction="row" className={classes.bgEven}>
            {renderOrders(t, doneOrders(retOrderList), 12, t('OrderStatus.DONE'))}
          </Grid>
          <Grid item container xs={3} spacing={1} alignItems="flex-start" alignContent="flex-start"
                justify="flex-start" direction="row" className={classes.bgOdd}>
            {renderOrders(t, exceptOrders(retOrderList), 12, t('OrderStatus.EXCP'))}
          </Grid>
        </Grid>
      </div>)}
  </I18n>;

}

const mapState = (state, props) => ({
  ...props,
  orderList: state.order.list
});

const mapDispatch = {
  getDetail: orderActions.getDetail,
  view: orderActions.view,
  doPush: push,
  getList: orderActions.getList
};

export default connect(mapState, mapDispatch)(HomeOperationList);
