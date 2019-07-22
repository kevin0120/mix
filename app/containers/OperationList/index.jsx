import React from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardMedia from '@material-ui/core/CardMedia';
import Typography from '@material-ui/core/Typography';
import { makeStyles, createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles';
import clsx from 'clsx';
import CardContent from '@material-ui/core/CardContent';
import { orderActions } from '../../modules/order/action';
import { ORDER_STATUS } from '../../modules/order/model';
import { todoOrders, doneOrders, excpOrders } from '../../modules/order/selector';
import styles from './styles';
import settingImg from '../../../resources/imgs/setting.png';

type Props = {
  classes: {},
  orderList: [],
  trigger: ()=>{},
  doPush: ()=>{}
};

function HomeOperationList(props: Props) {
  const classes = makeStyles(styles)();
  const { trigger, doPush, orderList } = props;

  const onCardClick = (order) => {
    trigger(order);
    doPush('/app/working');
  };

  const statusMap = {
    [ORDER_STATUS.TODO]: classes.statusTodo,
    [ORDER_STATUS.WIP]: classes.statusWIP,
    [ORDER_STATUS.DONE]: classes.statusDone,
    [ORDER_STATUS.CANCEL]: classes.statusCancel,
    [ORDER_STATUS.PENDING]: classes.statusPending,
    [ORDER_STATUS.FAIL]: classes.statusFail,
  };

  const renderOrders = (orders, size, title) => {
    return (
      <React.Fragment>
        <Grid item xs={12} className={classes.listTitle}>
          <Typography gutterBottom variant="h6" align="left">
            {title}
          </Typography>
        </Grid>
        {orders?.map((order) => <Grid item xs={size} key={order.name}>
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
                    {order.status || ORDER_STATUS.TODO}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Paper>
          </Grid>
        ) || null}
      </React.Fragment>
    );
  };

  return <div className={classes.root}>
    <Grid container className={clsx(classes.container, classes.bgEven)} justify="center" spacing={4}>
      <Grid item container xs={6} spacing={1} alignItems="flex-start" alignContent="flex-start"
            justify="flex-start" direction="row" className={classes.bgOdd}>
        {renderOrders(todoOrders(orderList), 6, 'TODO')}
      </Grid>
      <Grid item container xs={3} spacing={1} alignItems="flex-start" alignContent="flex-start"
            justify="flex-start" direction="row" className={classes.bgEven}>
        {renderOrders(doneOrders(orderList), 12, 'DONE')}
      </Grid>
      <Grid item container xs={3} spacing={1} alignItems="flex-start" alignContent="flex-start"
            justify="flex-start" direction="row" className={classes.bgOdd}>
        {renderOrders(excpOrders(orderList), 12, 'EXCP')}
      </Grid>
    </Grid>
  </div>;

}

const mapState = (state, props) => ({
  ...props,
  orderList: state.order.list
});

const mapDispatch = {
  trigger: orderActions.trigger,
  doPush: push
};

export default connect(mapState, mapDispatch)(HomeOperationList);
