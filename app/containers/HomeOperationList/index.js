import React from 'react';
import { connect } from 'react-redux';
import { orderActions } from '../../modules/order/action';
import { push } from 'connected-react-router';

class HomeOperationList extends React.Component {
  render() {
    const { orderList, trigger, doPush } = this.props;
    return <div>
      {orderList.map((order) => {
        return <button onClick={() => {
          trigger(order);
          doPush('/app/working');
        }}>{order.name}</button>;
      })}
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

export default connect(mapState, mapDispatch)(HomeOperationList);
