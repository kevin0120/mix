import React from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import Typography from '@material-ui/core/Typography';
import { orderActions } from '../../modules/order/action';
import { ORDER_STATUS } from '../../modules/order/model';
import styles from './styles';

type Props = {
  classes: {},
  orderList: [],
  trigger: ()=>{},
  doPush: ()=>{}
};

class HomeOperationList extends React.Component<Props> {
  onCardClick = (order) => {
    const { trigger, doPush } = this.props;
    trigger(order);
    doPush('/app/working');
  };

  renderOrders = (orders, size) => {
    const { classes } = this.props;
    return orders.map((order) => <Grid item xs={size} key={order.name} >
      <Card>
        <CardActionArea className={classes.orderCard} onClick={() => this.onCardClick(order)}>
          <Typography gutterBottom variant="h5" component="h2">
            {order.name}
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
            {order.info}
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
            {order.status || ORDER_STATUS.READY}
          </Typography>
        </CardActionArea>
      </Card>
    </Grid>);
  };

  render() {
    const { classes, orderList } = this.props;
    return <div className={classes.root}>
      <Grid container className={classes.container} justify="center" spacing={4}>
        <Grid item container xs={6} spacing={1} alignItems={'flex-start'} alignContent={'flex-start'} justify={'flex-start' } direction={'row'}>
          {this.renderOrders(orderList, 6)}
        </Grid>
        <Grid item container xs={3} spacing={1} alignItems={'flex-start'} alignContent={'flex-start'} justify={'flex-start'} direction={'row'}>
          {this.renderOrders(orderList, 12)}

        </Grid>
        <Grid item container xs={3} spacing={1} alignItems={'flex-start'} alignContent={'flex-start'} justify={'flex-start'} direction={'row'}>
          {this.renderOrders(orderList, 12)}
        </Grid>
      </Grid>
    </div>;
  }

}

const mapState = (state, props) => ({
  ...props,
  orderList: state.order.list
});

const mapDispatch = {
  trigger: orderActions.trigger,
  doPush: push
};

export default withStyles(styles)(connect(mapState, mapDispatch)(HomeOperationList));
