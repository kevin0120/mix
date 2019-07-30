import { connect } from 'react-redux';
import React from 'react';
import orderTypes from './orderTypes';

function ConnectedWorkingTemplate({ order }) {
  const orderType = order.viewingOrder && order.viewingOrder.type || 'default';
  const WorkingPage = (orderTypes[orderType] && orderTypes[orderType].container);
  return <WorkingPage/>;
}

const mapState = (state, props) => ({
  order: state.order,
  ...props
});

const mapDispatch = {};

export default connect(mapState, mapDispatch)(ConnectedWorkingTemplate);
