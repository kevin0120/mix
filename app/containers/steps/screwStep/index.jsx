import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import styles from './style';
import { viewingStep, stepPayload } from '../../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  points: stepPayload(viewingStep(state.order))?.points || [],
  image: stepPayload(viewingStep(state.order))?.image || {}
});

const mapDispatch = {};

function ScrewStep() {
  const classes = makeStyles(styles.layout)();

  return <div className={classes}>
    screwStep
  </div>;
}


export default connect(
  mapState,
  mapDispatch
)(ScrewStep);

