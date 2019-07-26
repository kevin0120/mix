import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import styles from './style';
import { viewingStep, stepPayload } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';

const mapState = (state, props) => ({
  ...props,
  points: stepPayload(viewingStep(state.order))?.points || [],
  image: stepPayload(viewingStep(state.order))?.image || {}
});

const mapDispatch = {};

function ScrewStep({ image, points }) {
  const classes = makeStyles(styles)();

  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
    />
  </div>;
}


export default connect(
  mapState,
  mapDispatch
)(ScrewStep);

