import { connect } from 'react-redux';
import React from 'react';
import orderTypes from './orderTypes';
import StepWorking from '../stepWorking';

class ConnectedWorkingTemplate extends React.Component {

  render() {
    const { order } = this.props;
    if (order.currentOrder && order.currentOrder.type) {
      const WorkingPage = orderTypes[order.currentOrder.type] && orderTypes[order.currentOrder.type].container || StepWorking;
      return <WorkingPage/>;
    } else {
      return null;
    }
  }
}

const mapState = (state, props) => ({
  order: state.order,
  ...props
});

const mapDispatch = {};

export default connect(mapState, mapDispatch)(ConnectedWorkingTemplate);
